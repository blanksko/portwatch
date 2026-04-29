// Package epoch tracks a monotonically incrementing scan counter per host.
// Each time a scan cycle completes for a host the epoch is bumped, giving
// downstream components a cheap way to detect whether results are stale.
package epoch

import (
	"sync"
)

// Tracker maintains per-host epoch counters.
type Tracker struct {
	mu      sync.Mutex
	counters map[string]uint64
}

// New returns an initialised Tracker.
func New() *Tracker {
	return &Tracker{
		counters: make(map[string]uint64),
	}
}

// Bump increments the epoch for host and returns the new value.
func (t *Tracker) Bump(host string) uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.counters[host]++
	return t.counters[host]
}

// Get returns the current epoch for host.
// If the host has never been seen, 0 is returned.
func (t *Tracker) Get(host string) uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.counters[host]
}

// Reset sets the epoch for host back to zero.
func (t *Tracker) Reset(host string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.counters, host)
}

// Snapshot returns a copy of all current counters.
func (t *Tracker) Snapshot() map[string]uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make(map[string]uint64, len(t.counters))
	for k, v := range t.counters {
		out[k] = v
	}
	return out
}

// Stale reports whether the epoch stored for host is older than the
// provided epoch value (i.e. the caller holds a result from a previous cycle).
func (t *Tracker) Stale(host string, epoch uint64) bool {
	return t.Get(host) > epoch
}
