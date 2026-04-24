// Package jitter provides randomised delay utilities to spread
// concurrent scan cycles and avoid thundering-herd effects when
// many hosts are polled on the same schedule.
package jitter

import (
	"math/rand"
	"sync"
	"time"
)

// Jitter adds a bounded random offset to a base duration.
type Jitter struct {
	mu      sync.Mutex
	rng     *rand.Rand
	maxFrac float64 // fraction of base, e.g. 0.25 = ±25 %
}

// New returns a Jitter that spreads delays by up to maxFraction of
// the requested base duration. maxFraction is clamped to [0, 1].
func New(maxFraction float64) *Jitter {
	if maxFraction < 0 {
		maxFraction = 0
	}
	if maxFraction > 1 {
		maxFraction = 1
	}
	return &Jitter{
		rng:     rand.New(rand.NewSource(time.Now().UnixNano())), //nolint:gosec
		maxFrac: maxFraction,
	}
}

// Apply returns base ± a random offset bounded by maxFraction * base.
// The result is always >= 0.
func (j *Jitter) Apply(base time.Duration) time.Duration {
	if base <= 0 || j.maxFrac == 0 {
		return base
	}
	j.mu.Lock()
	f := j.rng.Float64() // [0, 1)
	j.mu.Unlock()

	// Map f to [-maxFrac, +maxFrac]
	offset := time.Duration(float64(base) * j.maxFrac * (2*f - 1))
	result := base + offset
	if result < 0 {
		return 0
	}
	return result
}

// Sleep blocks for Apply(base).
func (j *Jitter) Sleep(base time.Duration) {
	time.Sleep(j.Apply(base))
}
