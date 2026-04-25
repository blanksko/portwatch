// Package decay provides a score-based host health tracker that decays
// scores over time, allowing portwatch to deprioritise hosts that have
// been consistently stable and re-elevate those that show new activity.
package decay

import (
	"sync"
	"time"
)

// DefaultDecayRate is the fraction by which a score is multiplied each
// decay interval when no new events are recorded.
const DefaultDecayRate = 0.5

// entry holds the current score and the last time it was updated.
type entry struct {
	score     float64
	updatedAt time.Time
}

// Tracker maintains per-host activity scores that decay exponentially
// between updates.
type Tracker struct {
	mu        sync.Mutex
	scores    map[string]entry
	rate      float64
	interval  time.Duration
	now       func() time.Time
}

// New returns a Tracker that applies decayRate per interval.
// decayRate must be in (0, 1]; values outside that range are clamped.
func New(decayRate float64, interval time.Duration) *Tracker {
	if decayRate <= 0 || decayRate > 1 {
		decayRate = DefaultDecayRate
	}
	if interval <= 0 {
		interval = time.Minute
	}
	return &Tracker{
		scores:   make(map[string]entry),
		rate:     decayRate,
		interval: interval,
		now:      time.Now,
	}
}

// Record adds delta to the host's current (decayed) score.
func (t *Tracker) Record(host string, delta float64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := t.now()
	e := t.decay(host, now)
	e.score += delta
	e.updatedAt = now
	t.scores[host] = e
}

// Score returns the current decayed score for host.
func (t *Tracker) Score(host string) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.decay(host, t.now()).score
}

// Reset removes the host's score entirely.
func (t *Tracker) Reset(host string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.scores, host)
}

// decay computes the exponentially decayed entry for host at time now.
// The caller must hold t.mu.
func (t *Tracker) decay(host string, now time.Time) entry {
	e, ok := t.scores[host]
	if !ok {
		return entry{updatedAt: now}
	}
	elapsed := now.Sub(e.updatedAt)
	if elapsed <= 0 {
		return e
	}
	periods := float64(elapsed) / float64(t.interval)
	// score *= rate^periods
	factor := 1.0
	for i := 0.0; i < periods; i++ {
		factor *= t.rate
	}
	e.score *= factor
	return e
}
