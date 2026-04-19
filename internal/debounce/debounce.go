// Package debounce prevents repeated notifications for the same event
// within a configurable quiet period.
package debounce

import (
	"sync"
	"time"
)

// Debouncer tracks recent events and suppresses duplicates within a window.
type Debouncer struct {
	mu     sync.Mutex
	window time.Duration
	seen   map[string]time.Time
	now    func() time.Time
}

// New returns a Debouncer with the given quiet window.
func New(window time.Duration) *Debouncer {
	return &Debouncer{
		window: window,
		seen:   make(map[string]time.Time),
		now:    time.Now,
	}
}

// Allow returns true if the key has not been seen within the quiet window.
// If allowed, the key's timestamp is updated.
func (d *Debouncer) Allow(key string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	if last, ok := d.seen[key]; ok && now.Sub(last) < d.window {
		return false
	}
	d.seen[key] = now
	return true
}

// Reset removes a key so the next call to Allow will succeed.
func (d *Debouncer) Reset(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.seen, key)
}

// Purge removes all keys whose last-seen time is older than the window.
func (d *Debouncer) Purge() {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.now()
	for k, t := range d.seen {
		if now.Sub(t) >= d.window {
			delete(d.seen, k)
		}
	}
}
