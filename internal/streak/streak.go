// Package streak tracks consecutive scan cycles in which a host
// has reported no port changes. A long unbroken streak can be used
// to reduce scan frequency or suppress low-priority alerts.
package streak

import (
	"sync"
	"time"
)

// Streak holds per-host consecutive-stable-cycle counters.
type Streak struct {
	mu      sync.Mutex
	counts  map[string]int
	updated map[string]time.Time
	now     func() time.Time
}

// New returns an initialised Streak tracker.
func New() *Streak {
	return &Streak{
		counts:  make(map[string]int),
		updated: make(map[string]time.Time),
		now:     time.Now,
	}
}

// Record increments the stable-cycle counter for host when changed is
// false, and resets it to zero when a change is detected.
func (s *Streak) Record(host string, changed bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if changed {
		s.counts[host] = 0
	} else {
		s.counts[host]++
	}
	s.updated[host] = s.now()
}

// Count returns the current consecutive-stable count for host.
// An unknown host returns 0.
func (s *Streak) Count(host string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.counts[host]
}

// LastSeen returns the time of the most recent Record call for host
// and whether the host has been seen at all.
func (s *Streak) LastSeen(host string) (time.Time, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.updated[host]
	return t, ok
}

// Reset clears the streak counter and last-seen time for host.
func (s *Streak) Reset(host string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.counts, host)
	delete(s.updated, host)
}

// All returns a snapshot of all current streak counts keyed by host.
func (s *Streak) All() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]int, len(s.counts))
	for k, v := range s.counts {
		out[k] = v
	}
	return out
}
