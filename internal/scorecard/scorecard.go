// Package scorecard tracks a reliability score for each scanned host,
// increasing on successful scans and decaying on consecutive failures.
package scorecard

import (
	"sync"
	"time"
)

const (
	defaultMax   = 100.0
	defaultDecay = 10.0
	defaultGain  = 5.0
)

// Score holds the current reliability value and metadata for a host.
type Score struct {
	Value     float64
	UpdatedAt time.Time
}

// Scorecard maintains per-host reliability scores.
type Scorecard struct {
	mu     sync.RWMutex
	scores map[string]Score
	max    float64
	decay  float64
	gain   float64
	now    func() time.Time
}

// Option is a functional option for Scorecard.
type Option func(*Scorecard)

// WithMax sets the maximum achievable score (default 100).
func WithMax(max float64) Option {
	return func(s *Scorecard) { s.max = max }
}

// WithDecay sets the penalty applied on a failed scan (default 10).
func WithDecay(d float64) Option {
	return func(s *Scorecard) { s.decay = d }
}

// WithGain sets the reward applied on a successful scan (default 5).
func WithGain(g float64) Option {
	return func(s *Scorecard) { s.gain = g }
}

// New creates a Scorecard with optional configuration.
func New(opts ...Option) *Scorecard {
	s := &Scorecard{
		scores: make(map[string]Score),
		max:    defaultMax,
		decay:  defaultDecay,
		gain:   defaultGain,
		now:    time.Now,
	}
	for _, o := range opts {
		o(s)
	}
	return s
}

// RecordSuccess increases the score for host by the configured gain, capped at max.
func (s *Scorecard) RecordSuccess(host string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cur := s.scores[host]
	cur.Value = min(cur.Value+s.gain, s.max)
	cur.UpdatedAt = s.now()
	s.scores[host] = cur
}

// RecordFailure decreases the score for host by the configured decay, floored at 0.
func (s *Scorecard) RecordFailure(host string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	cur := s.scores[host]
	cur.Value = max(cur.Value-s.decay, 0)
	cur.UpdatedAt = s.now()
	s.scores[host] = cur
}

// Get returns the current Score for host.
func (s *Scorecard) Get(host string) Score {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.scores[host]
}

// Reset removes the score entry for host.
func (s *Scorecard) Reset(host string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.scores, host)
}

// All returns a snapshot of all current scores.
func (s *Scorecard) All() map[string]Score {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]Score, len(s.scores))
	for k, v := range s.scores {
		out[k] = v
	}
	return out
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
