package core

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/MetaLoan/pp/internal/config"
	"github.com/MetaLoan/pp/internal/market"
	"github.com/gorilla/websocket"
)

// MarketData represents a single price tick
type MarketData struct {
	Symbol    string  `json:"symbol"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

// MarketEngine manages market data and subscribers
type MarketEngine struct {
	subscribers  map[*websocket.Conn]bool
	broadcast    chan []byte
	register     chan *websocket.Conn
	unregister   chan *websocket.Conn
	mu           sync.RWMutex
	lastPrice    map[string]MarketData
	provider     market.PriceProvider
	pollInterval time.Duration
	priceHistory map[string][]MarketData
}

var GlobalMarket *MarketEngine

func InitMarket() {
	cfg := config.GlobalConfig.Market
	poll := time.Second
	if d, err := time.ParseDuration(cfg.PollInterval); err == nil && d > 0 {
		poll = d
	}

	var primary market.PriceProvider
	timeout := 3 * time.Second
	if t, err := time.ParseDuration(cfg.External.Timeout); err == nil && t > 0 {
		timeout = t
	}
	switch strings.ToLower(cfg.External.Provider) {
	case "binance":
		path := cfg.External.Path
		if path == "" || path == "/price" {
			path = "/api/v3/ticker/price"
		}
		primary = market.NewBinanceProvider(cfg.External.BaseURL, path, timeout)
	case "twelvedata":
		apiKey := cfg.External.APIKey
		if apiKey == "" {
			// Also allow env-specific override commonly used
			if val := strings.TrimSpace(getEnvFallback("TWELVEDATA_API_KEY", "")); val != "" {
				apiKey = val
			}
		}
		primary = market.NewTwelveDataProvider(cfg.External.BaseURL, cfg.External.Path, apiKey, timeout)
	default:
		if cfg.External.BaseURL != "" {
			primary = market.NewHTTPProvider(cfg.External.BaseURL, cfg.External.Path, timeout)
		}
	}
	mockInit := cfg.MockInitial
	if mockInit == nil {
		mockInit = map[string]float64{}
	}
	// Ensure all allowed symbols have a seed price
	symbols := config.GlobalConfig.Trading.AllowedSymbols
	if len(symbols) == 0 {
		symbols = []string{"EURUSD"}
	}
	for _, sym := range symbols {
		if _, ok := mockInit[sym]; !ok {
			mockInit[sym] = market.SeedForSymbol(sym)
		}
	}
	cb := market.NewCircuitBreaker(cfg.Breaker.FailThreshold, cfg.Breaker.OpenSeconds)
	var provider market.PriceProvider
	switch {
	case primary != nil:
		// Use only the real provider; on error we'll reuse last price instead of random mock.
		provider = &market.CircuitGuard{Provider: primary, Breaker: cb}
	default:
		provider = market.NewMockProvider(mockInit)
	}

	GlobalMarket = &MarketEngine{
		subscribers:  make(map[*websocket.Conn]bool),
		broadcast:    make(chan []byte),
		register:     make(chan *websocket.Conn),
		unregister:   make(chan *websocket.Conn),
		lastPrice:    make(map[string]MarketData),
		provider:     provider,
		pollInterval: poll,
		priceHistory: make(map[string][]MarketData),
	}
	go GlobalMarket.run()
	go GlobalMarket.pollPrices()
}

func (m *MarketEngine) run() {
	for {
		select {
		case conn := <-m.register:
			m.mu.Lock()
			m.subscribers[conn] = true
			m.mu.Unlock()
		case conn := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.subscribers[conn]; ok {
				delete(m.subscribers, conn)
				conn.Close()
			}
			m.mu.Unlock()
		case message := <-m.broadcast:
			m.mu.RLock()
			for conn := range m.subscribers {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Websocket write error: %v", err)
					conn.Close()
					delete(m.subscribers, conn)
				}
			}
			m.mu.RUnlock()
		}
	}
}

func (m *MarketEngine) Register(conn *websocket.Conn) {
	m.register <- conn
}

func (m *MarketEngine) Unregister(conn *websocket.Conn) {
	m.unregister <- conn
}

func (m *MarketEngine) pollPrices() {
	// Use jittered sleeps to avoid syncing with upstream rate limits exactly.
	lastTS := map[string]int64{}
	for {
		symbols := config.GlobalConfig.Trading.AllowedSymbols
		if len(symbols) == 0 {
			symbols = []string{"EURUSD"}
		}
		for _, sym := range symbols {
			apiSymbol := sym
			// Support case-insensitive symbol map keys because Viper lowercases map keys by default.
			if mapped, ok := config.GlobalConfig.Market.SymbolMap[sym]; ok && mapped != "" {
				apiSymbol = mapped
			} else if mapped, ok := config.GlobalConfig.Market.SymbolMap[strings.ToLower(sym)]; ok && mapped != "" {
				apiSymbol = mapped
			}
			price, err := m.provider.GetPrice(apiSymbol)
			if err != nil {
				log.Printf("price fetch error for %s: %v; keeping previous value", sym, err)
				m.mu.RLock()
				last, ok := m.lastPrice[sym]
				m.mu.RUnlock()
				if !ok {
					// No previous price — skip this tick to avoid injecting a bogus seed.
					continue
				}
				price = last.Price
			}
			now := time.Now().UnixMilli()
			if last, ok := lastTS[sym]; ok && now <= last {
				now = last + 1
			}
			lastTS[sym] = now

			data := MarketData{
				Symbol:    sym,
				Price:     price,
				Timestamp: now,
			}

			m.mu.Lock()
			m.lastPrice[data.Symbol] = data
			m.priceHistory[data.Symbol] = append(m.priceHistory[data.Symbol], data)
			if len(m.priceHistory[data.Symbol]) > 5000 {
				m.priceHistory[data.Symbol] = m.priceHistory[data.Symbol][len(m.priceHistory[data.Symbol])-4000:]
			}
			m.mu.Unlock()

			msg, err := json.Marshal(data)
			if err != nil {
				log.Println("Error marshaling market data:", err)
				continue
			}

			m.broadcast <- msg
		}

		// Sleep with small jitter (±20%) to spread out API calls.
		pollMs := m.pollInterval.Milliseconds()
		jitter := pollMs / 5 // 20%
		if jitter < 1 {
			jitter = 1
		}
		sleepMs := pollMs + rand.Int63n(jitter*2) - jitter
		if sleepMs < 200 {
			sleepMs = 200
		}
		time.Sleep(time.Duration(sleepMs) * time.Millisecond)
	}
}

// GetCurrentPrice returns the latest price for a symbol
func (m *MarketEngine) GetCurrentPrice(symbol string) (float64, error) {
	m.mu.RLock()
	data, ok := m.lastPrice[symbol]
	m.mu.RUnlock()
	if !ok || data.Price == 0 {
		return 0, fmt.Errorf("no price available for %s", symbol)
	}
	maxBehind := time.Duration(config.GlobalConfig.Market.MaxBehindSeconds) * time.Second
	if maxBehind > 0 && time.Since(time.UnixMilli(data.Timestamp)) > maxBehind {
		return 0, fmt.Errorf("stale price for %s", symbol)
	}
	return data.Price, nil
}

// GetCandles aggregates in-memory ticks into OHLC candles.
func (m *MarketEngine) GetCandles(symbol string, intervalSec int, limit int) []map[string]interface{} {
	if intervalSec <= 0 {
		intervalSec = 5
	}
	m.mu.RLock()
	history := append([]MarketData(nil), m.priceHistory[symbol]...)
	m.mu.RUnlock()
	if len(history) == 0 {
		return []map[string]interface{}{}
	}
	buckets := map[int64]*struct {
		open  float64
		high  float64
		low   float64
		close float64
		time  int64
	}{}
	for _, p := range history {
		bucket := (p.Timestamp / 1000 / int64(intervalSec)) * int64(intervalSec)
		if _, ok := buckets[bucket]; !ok {
			buckets[bucket] = &struct {
				open  float64
				high  float64
				low   float64
				close float64
				time  int64
			}{
				open:  p.Price,
				high:  p.Price,
				low:   p.Price,
				close: p.Price,
				time:  bucket,
			}
		} else {
			b := buckets[bucket]
			if p.Price > b.high {
				b.high = p.Price
			}
			if p.Price < b.low {
				b.low = p.Price
			}
			b.close = p.Price
		}
	}
	out := make([]map[string]interface{}, 0, len(buckets))
	for _, v := range buckets {
		out = append(out, map[string]interface{}{
			"time":  v.time,
			"open":  v.open,
			"high":  v.high,
			"low":   v.low,
			"close": v.close,
		})
	}
	// sort by time
	for i := 0; i < len(out)-1; i++ {
		for j := i + 1; j < len(out); j++ {
			if out[i]["time"].(int64) > out[j]["time"].(int64) {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	if limit > 0 && len(out) > limit {
		out = out[len(out)-limit:]
	}
	return out
}

func getEnvFallback(key, def string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return def
}
