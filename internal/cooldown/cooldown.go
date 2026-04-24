// Package cooldown provides per-host alert cooldown tracking to prevent
// repeated notifications for the same condition within a configurable window.
package cooldown

import (
	"sync"
	"time"
)

// Clock is a function that returns the current time.
type Clock func() time.Time

// Cooldown tracks the last alert time per host and suppresses repeated
// alerts until the cooldown window has elapsed.
type Cooldown struct {
	mu       sync.Mutex
	last     map[string]time.Time
	window   time.Duration
	now      Clock
}

// New returns a Cooldown with the given window duration.
// If window is zero or negative, a default of 5 minutes is used.
func New(window time.Duration, now Clock) *Cooldown {
	if window <= 0 {
		window = 5 * time.Minute
	}
	if now == nil {
		now = time.Now
	}
	return &Cooldown{
		last:   make(map[string]time.Time),
		window: window,
		now:    now,
	}
}

// Allow returns true if the host is not in cooldown and records the current
// time as the last alert time. Returns false if the cooldown window has not
// yet elapsed since the last alert.
func (c *Cooldown) Allow(host string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := c.now()
	if last, ok := c.last[host]; ok {
		if now.Sub(last) < c.window {
			return false
		}
	}
	c.last[host] = now
	return true
}

// Reset clears the cooldown state for the given host, allowing the next
// call to Allow to succeed immediately.
func (c *Cooldown) Reset(host string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.last, host)
}

// Remaining returns the duration remaining in the cooldown window for the
// given host. Returns zero if the host is not in cooldown.
func (c *Cooldown) Remaining(host string) time.Duration {
	c.mu.Lock()
	defer c.mu.Unlock()

	last, ok := c.last[host]
	if !ok {
		return 0
	}
	remaining := c.window - c.now().Sub(last)
	if remaining < 0 {
		return 0
	}
	return remaining
}
