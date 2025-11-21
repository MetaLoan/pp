package core

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

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
	subscribers map[*websocket.Conn]bool
	broadcast   chan []byte
	register    chan *websocket.Conn
	unregister  chan *websocket.Conn
	mu          sync.RWMutex
	lastPrice   map[string]MarketData
}

var GlobalMarket *MarketEngine

func InitMarket() {
	rand.Seed(time.Now().UnixNano())
	GlobalMarket = &MarketEngine{
		subscribers: make(map[*websocket.Conn]bool),
		broadcast:   make(chan []byte),
		register:    make(chan *websocket.Conn),
		unregister:  make(chan *websocket.Conn),
		lastPrice:   make(map[string]MarketData),
	}
	go GlobalMarket.run()
	go GlobalMarket.generateMockData()
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

// generateMockData simulates market data updates
func (m *MarketEngine) generateMockData() {
	// Simple multi-symbol random walk
	initialPrices := map[string]float64{
		"EURUSD":  1.0500,
		"BTCUSDT": 65000,
		"ETHUSDT": 3200,
	}

	for symbol, startPrice := range initialPrices {
		go func(sym string, price float64) {
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()
			for range ticker.C {
				// Random walk with symbol-specific volatility
				step := (rand.Float64() - 0.5)
				switch sym {
				case "EURUSD":
					price += step * 0.0010
				case "BTCUSDT":
					price += step * 50
				case "ETHUSDT":
					price += step * 2.5
				default:
					price += step * 0.1
				}

				data := MarketData{
					Symbol:    sym,
					Price:     price,
					Timestamp: time.Now().UnixMilli(),
				}

				m.mu.Lock()
				m.lastPrice[data.Symbol] = data
				m.mu.Unlock()

				msg, err := json.Marshal(data)
				if err != nil {
					log.Println("Error marshaling market data:", err)
					continue
				}

				m.broadcast <- msg
			}
		}(symbol, startPrice)
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
	return data.Price, nil
}
