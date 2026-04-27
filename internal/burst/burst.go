// Package burst provides a token-bucket style burst detector that flags
// hosts emitting an unusually high number of port-change events within a
// short sliding window.
package burst

import (
	"sync"
	"time"
)

// Detector tracks per-host event counts and reports when a host exceeds the
// configured burst threshold within the observation window.
type Detector struct {
	mu        sync.Mutex
	events    map[string][]time.Time
	threshold int
	window    time.Duration
	now       func() time.Time
}

// Option configures a Detector.
type Option func(*Detector)

// WithNow overrides the clock used by the Detector (useful in tests).
func WithNow(fn func() time.Time) Option {
	return func(d *Detector) { d.now = fn }
}

// New returns a Detector that fires when more than threshold events are
// recorded for the same host within window.
func New(threshold int, window time.Duration, opts ...Option) *Detector {
	if threshold <= 0 {
		threshold = 5
	}
	if window <= 0 {
		window = time.Minute
	}
	d := &Detector{
		events:    make(map[string][]time.Time),
		threshold: threshold,
		window:    window,
		now:       time.Now,
	}
	for _, o := range opts {
		o(d)
	}
	return d
}

// Record registers one event for host and returns true when the burst
// threshold has been exceeded within the current window.
func (d *Detector) Record(host string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	cutoff := now.Add(-d.window)

	ts := d.events[host]
	// prune expired entries
	valid := ts[:0]
	for _, t := range ts {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	valid = append(valid, now)
	d.events[host] = valid

	return len(valid) > d.threshold
}

// Count returns the number of events recorded for host within the current window.
func (d *Detector) Count(host string) int {
	d.mu.Lock()
	defer d.mu.Unlock()

	now := d.now()
	cutoff := now.Add(-d.window)
	count := 0
	for _, t := range d.events[host] {
		if t.After(cutoff) {
			count++
		}
	}
	return count
}

// Reset clears all recorded events for host.
func (d *Detector) Reset(host string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.events, host)
}
