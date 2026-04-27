// Package circuit provides a per-host circuit breaker that tracks scan
// success/failure rates and opens the circuit when a host becomes
// consistently unreachable, preventing wasted scan attempts.
package circuit

import (
	"sync"
	"time"
)

// State represents the current circuit state for a host.
type State int

const (
	StateClosed   State = iota // normal operation
	StateOpen                  // host is blocked
	StateHalfOpen              // probing after cooldown
)

// String returns a human-readable state name.
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Breaker is a per-host circuit breaker.
type Breaker struct {
	mu        sync.Mutex
	hosts     map[string]*entry
	threshold int
	cooldown  time.Duration
	now       func() time.Time
}

type entry struct {
	failures  int
	state     State
	openedAt  time.Time
}

// New creates a Breaker that opens after threshold consecutive failures
// and attempts recovery after cooldown.
func New(threshold int, cooldown time.Duration) *Breaker {
	return &Breaker{
		hosts:     make(map[string]*entry),
		threshold: threshold,
		cooldown:  cooldown,
		now:       time.Now,
	}
}

// Allow reports whether a scan attempt should proceed for host.
// It transitions the circuit to half-open when the cooldown has elapsed.
func (b *Breaker) Allow(host string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	e := b.get(host)
	switch e.state {
	case StateClosed, StateHalfOpen:
		return true
	case StateOpen:
		if b.now().Sub(e.openedAt) >= b.cooldown {
			e.state = StateHalfOpen
			return true
		}
		return false
	}
	return true
}

// RecordSuccess resets the failure counter and closes the circuit for host.
func (b *Breaker) RecordSuccess(host string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	e := b.get(host)
	e.failures = 0
	e.state = StateClosed
}

// RecordFailure increments the failure counter and opens the circuit if
// the threshold has been reached.
func (b *Breaker) RecordFailure(host string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	e := b.get(host)
	e.failures++
	if e.failures >= b.threshold && e.state != StateOpen {
		e.state = StateOpen
		e.openedAt = b.now()
	}
}

// StateOf returns the current circuit state for host.
func (b *Breaker) StateOf(host string) State {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.get(host).state
}

func (b *Breaker) get(host string) *entry {
	if e, ok := b.hosts[host]; ok {
		return e
	}
	e := &entry{state: StateClosed}
	b.hosts[host] = e
	return e
}
