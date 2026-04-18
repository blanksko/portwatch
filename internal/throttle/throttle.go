// Package throttle provides a scan throttle that limits how frequently
// a host can be scanned within a rolling time window.
package throttle

import (
	"sync"
	"time"
)

// Throttle tracks last-scan times per host and enforces a minimum interval.
type Throttle struct {
	mu       sync.Mutex
	last     map[string]time.Time
	interval time.Duration
	now      func() time.Time
}

// New returns a Throttle that enforces the given minimum interval between scans.
func New(interval time.Duration) *Throttle {
	return &Throttle{
		last:     make(map[string]time.Time),
		interval: interval,
		now:      time.Now,
	}
}

// Allow reports whether host may be scanned now.
// If allowed, the last-scan time is updated.
func (t *Throttle) Allow(host string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	if last, ok := t.last[host]; ok {
		if now.Sub(last) < t.interval {
			return false
		}
	}
	t.last[host] = now
	return true
}

// Reset clears the recorded scan time for host, allowing an immediate scan.
func (t *Throttle) Reset(host string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.last, host)
}

// NextAllowed returns the earliest time host may be scanned again.
// If the host has never been scanned, it returns the zero time.
func (t *Throttle) NextAllowed(host string) time.Time {
	t.mu.Lock()
	defer t.mu.Unlock()
	if last, ok := t.last[host]; ok {
		return last.Add(t.interval)
	}
	return time.Time{}
}
