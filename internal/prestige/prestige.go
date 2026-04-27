// Package prestige tracks the reputation score of a host based on how
// frequently unexpected port changes are observed. Hosts with repeated
// deviations accumulate a higher prestige penalty, which callers can use
// to prioritise alerting or throttle noisy sources.
package prestige

import (
	"sync"
	"time"
)

// Record holds the reputation state for a single host.
type Record struct {
	Score     float64
	Deviations int
	LastSeen  time.Time
}

// Tracker maintains reputation scores for monitored hosts.
type Tracker struct {
	mu      sync.RWMutex
	records map[string]*Record
	decay   float64 // fraction subtracted per interval
	gain    float64 // added per deviation event
	now     func() time.Time
}

// Option configures a Tracker.
type Option func(*Tracker)

// WithDecay sets the per-interval decay fraction (default 0.1).
func WithDecay(d float64) Option {
	return func(t *Tracker) { t.decay = d }
}

// WithGain sets the score added per deviation (default 1.0).
func WithGain(g float64) Option {
	return func(t *Tracker) { t.gain = g }
}

// WithNow overrides the clock (useful in tests).
func WithNow(fn func() time.Time) Option {
	return func(t *Tracker) { t.now = fn }
}

// New returns a ready-to-use Tracker.
func New(opts ...Option) *Tracker {
	t := &Tracker{
		records: make(map[string]*Record),
		decay:   0.1,
		gain:    1.0,
		now:     time.Now,
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Record registers a deviation event for host and returns the updated score.
func (t *Tracker) Record(host string) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	r, ok := t.records[host]
	if !ok {
		r = &Record{}
		t.records[host] = r
	}
	r.Score += t.gain
	r.Deviations++
	r.LastSeen = t.now()
	return r.Score
}

// Decay reduces the score of all hosts by the configured decay fraction.
// Call this periodically (e.g. once per scan cycle) to let quiet hosts
// recover their reputation over time.
func (t *Tracker) Decay() {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, r := range t.records {
		r.Score -= r.Score * t.decay
		if r.Score < 0 {
			r.Score = 0
		}
	}
}

// Get returns the current Record for host. The second return value is false
// when the host has never been seen.
func (t *Tracker) Get(host string) (Record, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	r, ok := t.records[host]
	if !ok {
		return Record{}, false
	}
	return *r, true
}

// Reset removes all state for host.
func (t *Tracker) Reset(host string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.records, host)
}
