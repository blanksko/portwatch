// Package window provides a sliding time-window counter for tracking
// per-host event frequency over a configurable duration.
package window

import (
	"sync"
	"time"
)

// entry holds a single timestamped event.
type entry struct {
	at time.Time
}

// Counter tracks events within a sliding time window.
type Counter struct {
	mu     sync.Mutex
	window time.Duration
	events map[string][]entry
	now    func() time.Time
}

// New returns a Counter with the given sliding window duration.
// Panics if window is zero or negative.
func New(window time.Duration) *Counter {
	if window <= 0 {
		panic("window: duration must be positive")
	}
	return &Counter{
		window: window,
		events: make(map[string][]entry),
		now:    time.Now,
	}
}

// Add records one event for the given key.
func (c *Counter) Add(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prune(key)
	c.events[key] = append(c.events[key], entry{at: c.now()})
}

// Count returns the number of events recorded for key within the window.
func (c *Counter) Count(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.prune(key)
	return len(c.events[key])
}

// Reset removes all recorded events for key.
func (c *Counter) Reset(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.events, key)
}

// prune removes events outside the current window for key.
// Caller must hold c.mu.
func (c *Counter) prune(key string) {
	cutoff := c.now().Add(-c.window)
	ev := c.events[key]
	i := 0
	for i < len(ev) && ev[i].at.Before(cutoff) {
		i++
	}
	if i > 0 {
		c.events[key] = ev[i:]
	}
}
