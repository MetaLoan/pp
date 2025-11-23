package market

import (
	"sync"
	"time"
)

type CircuitBreaker struct {
	mu            sync.RWMutex
	failures      int
	state         string // closed, open, half-open
	lastFailure   time.Time
	openUntil     time.Time
	failThreshold int
	openSeconds   int
}

func NewCircuitBreaker(failThreshold, openSeconds int) *CircuitBreaker {
	return &CircuitBreaker{
		state:         "closed",
		failThreshold: failThreshold,
		openSeconds:   openSeconds,
	}
}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	if cb.state == "open" && time.Now().Before(cb.openUntil) {
		return false
	}
	return true
}

func (cb *CircuitBreaker) Success() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures = 0
	cb.state = "closed"
}

func (cb *CircuitBreaker) Fail() {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.failures++
	cb.lastFailure = time.Now()
	if cb.failures >= cb.failThreshold {
		cb.state = "open"
		cb.openUntil = time.Now().Add(time.Duration(cb.openSeconds) * time.Second)
	}
}
