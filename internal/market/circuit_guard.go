package market

import "fmt"

// CircuitGuard wraps a provider with a circuit breaker, but does not alter data.
// Errors propagate to the caller, which can decide how to reuse last price.
type CircuitGuard struct {
	Provider PriceProvider
	Breaker  *CircuitBreaker
}

func (g *CircuitGuard) GetPrice(symbol string) (float64, error) {
	if g.Breaker != nil && !g.Breaker.Allow() {
		return 0, fmt.Errorf("circuit open")
	}
	price, err := g.Provider.GetPrice(symbol)
	if err == nil {
		if g.Breaker != nil {
			g.Breaker.Success()
		}
		return price, nil
	}
	if g.Breaker != nil {
		g.Breaker.Fail()
	}
	return 0, err
}
