// Package ratelimit provides a simple token-bucket rate limiter
// to throttle alert notifications and avoid alert fatigue.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter controls how frequently events are allowed through.
type Limiter struct {
	mu       sync.Mutex
	interval time.Duration
	last     map[string]time.Time
}

// New creates a Limiter that allows at most one event per key per interval.
func New(interval time.Duration) *Limiter {
	if interval <= 0 {
		interval = time.Minute
	}
	return &Limiter{
		interval: interval,
		last:     make(map[string]time.Time),
	}
}

// Allow reports whether the event identified by key should be allowed
// through based on the configured interval. It updates the last-seen
// timestamp when it returns true.
func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	if t, ok := l.last[key]; ok && now.Sub(t) < l.interval {
		return false
	}
	l.last[key] = now
	return true
}

// Reset clears the rate-limit state for a specific key.
func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, key)
}

// ResetAll clears all rate-limit state.
func (l *Limiter) ResetAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[string]time.Time)
}
