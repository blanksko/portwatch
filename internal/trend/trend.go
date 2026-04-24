// Package trend tracks port change frequency over time for a host,
// allowing callers to detect hosts with unusually volatile port sets.
package trend

import (
	"sync"
	"time"
)

// Entry records a single change event for a host.
type Entry struct {
	Host      string
	ChangedAt time.Time
	Opened    int
	Closed    int
}

// Tracker accumulates change events and exposes frequency queries.
type Tracker struct {
	mu      sync.Mutex
	events  map[string][]Entry
	window  time.Duration
	nowFunc func() time.Time
}

// New returns a Tracker that considers events within window as recent.
func New(window time.Duration) *Tracker {
	return &Tracker{
		events:  make(map[string][]Entry),
		window:  window,
		nowFunc: time.Now,
	}
}

// Record appends a change event for the given host.
func (t *Tracker) Record(host string, opened, closed int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.nowFunc()
	t.events[host] = append(t.events[host], Entry{
		Host:      host,
		ChangedAt: now,
		Opened:    opened,
		Closed:    closed,
	})
	t.prune(host, now)
}

// Count returns the number of change events recorded for host within the window.
func (t *Tracker) Count(host string) int {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.prune(host, t.nowFunc())
	return len(t.events[host])
}

// Recent returns a copy of all events within the window for host.
func (t *Tracker) Recent(host string) []Entry {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.prune(host, t.nowFunc())
	src := t.events[host]
	out := make([]Entry, len(src))
	copy(out, src)
	return out
}

// Reset clears all recorded events for host.
func (t *Tracker) Reset(host string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.events, host)
}

// prune removes events outside the window. Caller must hold t.mu.
func (t *Tracker) prune(host string, now time.Time) {
	cutoff := now.Add(-t.window)
	evs := t.events[host]
	i := 0
	for i < len(evs) && evs[i].ChangedAt.Before(cutoff) {
		i++
	}
	if i > 0 {
		t.events[host] = evs[i:]
	}
}
