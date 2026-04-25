// Package quota enforces per-host scan quotas within a rolling time window.
// It prevents any single host from consuming disproportionate scan resources.
package quota

import (
	"fmt"
	"sync"
	"time"
)

// Quota tracks scan counts per host and enforces a maximum within a window.
type Quota struct {
	mu      sync.Mutex
	events  map[string][]time.Time
	max     int
	window  time.Duration
	now     func() time.Time
}

// Option configures a Quota.
type Option func(*Quota)

// WithNow overrides the clock (useful for testing).
func WithNow(fn func() time.Time) Option {
	return func(q *Quota) { q.now = fn }
}

// New creates a Quota that allows at most max scans per host within window.
func New(max int, window time.Duration, opts ...Option) *Quota {
	if max <= 0 {
		max = 1
	}
	if window <= 0 {
		window = time.Minute
	}
	q := &Quota{
		events: make(map[string][]time.Time),
		max:    max,
		window: window,
		now:    time.Now,
	}
	for _, o := range opts {
		o(q)
	}
	return q
}

// Allow returns nil if the host is within quota, or an error if the limit has
// been reached for the current window.
func (q *Quota) Allow(host string) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := q.now()
	cutoff := now.Add(-q.window)

	evts := q.events[host]
	filtered := evts[:0]
	for _, t := range evts {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}

	if len(filtered) >= q.max {
		q.events[host] = filtered
		return fmt.Errorf("quota: host %q exceeded %d scans per %s", host, q.max, q.window)
	}

	q.events[host] = append(filtered, now)
	return nil
}

// Reset clears the quota counters for a specific host.
func (q *Quota) Reset(host string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.events, host)
}

// Count returns the number of scans recorded for host within the current window.
func (q *Quota) Count(host string) int {
	q.mu.Lock()
	defer q.mu.Unlock()

	now := q.now()
	cutoff := now.Add(-q.window)
	count := 0
	for _, t := range q.events[host] {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}
