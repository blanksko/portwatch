// Package circuitbreaker provides a simple circuit breaker to stop
// scanning hosts that repeatedly fail, reducing noise and load.
package circuitbreaker

import (
	"fmt"
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed State = iota // normal operation
	StateOpen               // failures exceeded threshold; calls rejected
)

// Breaker tracks per-host failure counts and opens the circuit when the
// failure threshold is exceeded within the observation window.
type Breaker struct {
	mu        sync.Mutex
	threshold int
	ttl       time.Duration
	now       func() time.Time
	entries   map[string]*entry
}

type entry struct {
	failures  int
	openedAt  time.Time
	state     State
}

// New returns a Breaker that opens after threshold consecutive failures and
// resets after ttl has elapsed.
func New(threshold int, ttl time.Duration) *Breaker {
	if threshold <= 0 {
		panic("circuitbreaker: threshold must be > 0")
	}
	return &Breaker{
		threshold: threshold,
		ttl:       ttl,
		now:       time.Now,
		entries:   make(map[string]*entry),
	}
}

// Allow returns true if the host is allowed to be scanned.
func (b *Breaker) Allow(host string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	e, ok := b.entries[host]
	if !ok {
		return true
	}
	if e.state == StateOpen {
		if b.now().Sub(e.openedAt) >= b.ttl {
			delete(b.entries, host)
			return true
		}
		return false
	}
	return true
}

// RecordSuccess resets the failure count for host.
func (b *Breaker) RecordSuccess(host string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.entries, host)
}

// RecordFailure increments the failure count and opens the circuit if the
// threshold is reached.
func (b *Breaker) RecordFailure(host string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	e, ok := b.entries[host]
	if !ok {
		e = &entry{}
		b.entries[host] = e
	}
	e.failures++
	if e.failures >= b.threshold {
		e.state = StateOpen
		e.openedAt = b.now()
	}
}

// Status returns a human-readable status string for host.
func (b *Breaker) Status(host string) string {
	b.mu.Lock()
	defer b.mu.Unlock()
	e, ok := b.entries[host]
	if !ok {
		return fmt.Sprintf("%s: closed (0 failures)", host)
	}
	return fmt.Sprintf("%s: state=%d failures=%d", host, e.state, e.failures)
}
