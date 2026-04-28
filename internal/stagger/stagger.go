// Package stagger spreads scan jobs across a time window to avoid
// thundering-herd bursts when many hosts are scheduled simultaneously.
package stagger

import (
	"sync"
	"time"
)

// Stagger distributes per-host delays evenly across a fixed window.
type Stagger struct {
	mu     sync.Mutex
	window time.Duration
	slots  map[string]time.Duration
	next   time.Duration
	step   time.Duration
}

// New creates a Stagger that spreads n hosts across the given window.
// If n is zero or negative it defaults to 1.
func New(window time.Duration, n int) *Stagger {
	if n <= 0 {
		n = 1
	}
	return &Stagger{
		window: window,
		slots:  make(map[string]time.Duration),
		step:   window / time.Duration(n),
	}
}

// Delay returns the delay that should be applied before scanning host.
// Each unique host receives a stable, deterministic offset within the window.
// Subsequent calls for the same host return the same value.
func (s *Stagger) Delay(host string) time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()

	if d, ok := s.slots[host]; ok {
		return d
	}

	d := s.next
	s.slots[host] = d
	s.next += s.step
	if s.next >= s.window {
		s.next = 0
	}
	return d
}

// Reset clears all assigned slots so hosts are re-distributed on the next call.
func (s *Stagger) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.slots = make(map[string]time.Duration)
	s.next = 0
}

// Len returns the number of hosts that have been assigned a slot.
func (s *Stagger) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.slots)
}
