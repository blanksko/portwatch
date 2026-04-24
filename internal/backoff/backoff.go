// Package backoff provides exponential backoff with jitter for
// controlling retry intervals during repeated scan failures.
package backoff

import (
	"math"
	"sync"
	"time"
)

const (
	defaultBase    = 1 * time.Second
	defaultMax     = 60 * time.Second
	defaultFactor  = 2.0
	defaultJitter  = 0.2
)

// Backoff tracks per-host exponential backoff state.
type Backoff struct {
	mu      sync.Mutex
	base    time.Duration
	max     time.Duration
	factor  float64
	jitter  float64
	counts  map[string]int
	now     func() time.Time
}

// New returns a Backoff with sensible defaults.
func New(base, max time.Duration, factor, jitter float64) *Backoff {
	if base <= 0 {
		base = defaultBase
	}
	if max <= 0 {
		max = defaultMax
	}
	if factor < 1 {
		factor = defaultFactor
	}
	if jitter < 0 || jitter > 1 {
		jitter = defaultJitter
	}
	return &Backoff{
		base:   base,
		max:    max,
		factor: factor,
		jitter: jitter,
		counts: make(map[string]int),
		now:    time.Now,
	}
}

// Default returns a Backoff using package-level defaults.
func Default() *Backoff {
	return New(defaultBase, defaultMax, defaultFactor, defaultJitter)
}

// Next returns the next backoff duration for the given host and
// increments the failure counter.
func (b *Backoff) Next(host string) time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()

	n := b.counts[host]
	b.counts[host] = n + 1

	raw := float64(b.base) * math.Pow(b.factor, float64(n))
	if raw > float64(b.max) {
		raw = float64(b.max)
	}

	// apply ±jitter
	spread := raw * b.jitter
	// deterministic-ish: use nanosecond remainder as cheap pseudo-random
	nano := float64(b.now().UnixNano() % 1000)
	offset := spread * (nano/999.0*2 - 1)

	d := time.Duration(raw + offset)
	if d < 0 {
		d = b.base
	}
	return d
}

// Reset clears the failure counter for the given host.
func (b *Backoff) Reset(host string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.counts, host)
}

// Failures returns the current failure count for host.
func (b *Backoff) Failures(host string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.counts[host]
}
