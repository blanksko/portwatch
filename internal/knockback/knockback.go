// Package knockback implements an exponential back-off gate that temporarily
// blocks a host after repeated scan failures, preventing thundering-herd
// behaviour when a target becomes unreachable.
package knockback

import (
	"sync"
	"time"
)

// DefaultThreshold is the number of consecutive failures before a host is
// blocked.
const DefaultThreshold = 3

// DefaultBase is the initial back-off duration.
const DefaultBase = 30 * time.Second

// DefaultMax is the ceiling for the exponential back-off.
const DefaultMax = 10 * time.Minute

// Gate tracks consecutive failures per host and decides whether a scan should
// be permitted.
type Gate struct {
	mu        sync.Mutex
	threshold int
	base      time.Duration
	max       time.Duration
	state     map[string]*entry
	now       func() time.Time
}

type entry struct {
	failures  int
	blockedAt time.Time
	backoff   time.Duration
}

// New returns a Gate with default parameters.
func New() *Gate {
	return &Gate{
		threshold: DefaultThreshold,
		base:      DefaultBase,
		max:       DefaultMax,
		state:     make(map[string]*entry),
		now:       time.Now,
	}
}

// Allow returns true when the host is permitted to be scanned.  It must be
// paired with a call to RecordSuccess or RecordFailure after the scan
// completes.
func (g *Gate) Allow(host string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()
	e, ok := g.state[host]
	if !ok {
		return true
	}
	if e.failures < g.threshold {
		return true
	}
	return g.now().After(e.blockedAt.Add(e.backoff))
}

// RecordFailure increments the failure counter for host and, if the threshold
// is reached, sets the block window using exponential back-off.
func (g *Gate) RecordFailure(host string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	e, ok := g.state[host]
	if !ok {
		e = &entry{}
		g.state[host] = e
	}
	e.failures++
	if e.failures >= g.threshold {
		if e.backoff == 0 {
			e.backoff = g.base
		} else {
			e.backoff *= 2
			if e.backoff > g.max {
				e.backoff = g.max
			}
		}
		e.blockedAt = g.now()
	}
}

// RecordSuccess resets the failure state for host.
func (g *Gate) RecordSuccess(host string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.state, host)
}

// Failures returns the current consecutive failure count for host.
func (g *Gate) Failures(host string) int {
	g.mu.Lock()
	defer g.mu.Unlock()
	if e, ok := g.state[host]; ok {
		return e.failures
	}
	return 0
}
