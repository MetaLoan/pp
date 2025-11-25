package market

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

// MockProvider updates prices via a random walk; used as fallback or dev provider.
type MockProvider struct {
	mu     sync.RWMutex
	prices map[string]float64
}

func NewMockProvider(initial map[string]float64) *MockProvider {
	mp := &MockProvider{
		prices: map[string]float64{},
	}
	for k, v := range initial {
		mp.prices[k] = v
	}
	go mp.tick()
	return mp
}

func (m *MockProvider) tick() {
	// Reduce tick frequency to 1s to lower simulated update rate
	// Tick every 200ms as requested (0.2s)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()
	for range ticker.C {
		m.mu.Lock()
		for sym, price := range m.prices {
			// Apply a small random percentage change +/-0.005% per tick
			// 0.005% = 0.00005 in decimal
			pct := (rand.Float64()*2 - 1) * 0.00005
			price += price * pct
			m.prices[sym] = price
			// Debug: log EURUSD updates to trace unexpected large jumps
			if sym == "EURUSD" || sym == "EUR/USD" {
				log.Printf("[mock] tick %s -> %.8f (pct %.8f)", sym, price, pct)
			}
		}
		m.mu.Unlock()
	}
}

func (m *MockProvider) GetPrice(symbol string) (float64, error) {
	m.mu.RLock()
	if p, ok := m.prices[symbol]; ok {
		m.mu.RUnlock()
		return p, nil
	}
	m.mu.RUnlock()

	// Seed a default price to avoid hard failures and let tick mutate it.
	p := SeedForSymbol(symbol)
	m.mu.Lock()
	m.prices[symbol] = p
	m.mu.Unlock()
	return p, nil
}

func SeedForSymbol(symbol string) float64 {
	// Use 10000 as the seed for all symbols per request
	return 10000
}
