package market

import (
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
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		m.mu.Lock()
		for sym, price := range m.prices {
			// Increased volatility for faster price movement simulation
			step := (rand.Float64() - 0.5)
			switch sym {
			case "EURUSD", "EUR/USD":
				price += step * 0.0005
			case "BTCUSDT", "BTC/USD":
				price += step * 50.0
			case "ETHUSDT", "ETH/USD":
				price += step * 5.0
			case "S1":
				price += step * 2.0
			default:
				price += step * 0.01 // Reduced default volatility
			}
			m.prices[sym] = price
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
	switch symbol {
	case "BTCUSDT", "BTC/USD", "BTCUSDC":
		return 65000 + (rand.Float64()-0.5)*500
	case "ETHUSDT", "ETH/USD", "ETHUSDC":
		return 3200 + (rand.Float64()-0.5)*50
	case "S1":
		return 100 + (rand.Float64()-0.5)*10
	case "EURUSD":
		return 1.05 + (rand.Float64()-0.5)*0.01
	default:
		return 1.0 + rand.Float64()*0.05
	}
}
