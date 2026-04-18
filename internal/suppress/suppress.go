// Package suppress provides a mechanism to suppress repeated alerts
// for the same host within a configurable cooldown window.
package suppress

import (
	"sync"
	"time"
)

// Suppressor tracks alert state per host and suppresses duplicate
// notifications until the cooldown period has elapsed.
type Suppressor struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
	now      func() time.Time
}

// New returns a Suppressor with the given cooldown duration.
func New(cooldown time.Duration) *Suppressor {
	return &Suppressor{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if an alert for the given host should be sent.
// It returns false if the host was alerted within the cooldown window.
func (s *Suppressor) Allow(host string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, seen := s.last[host]
	if seen && s.now().Sub(t) < s.cooldown {
		return false
	}
	s.last[host] = s.now()
	return true
}

// Reset clears the suppression state for the given host.
func (s *Suppressor) Reset(host string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.last, host)
}

// ResetAll clears suppression state for all hosts.
func (s *Suppressor) ResetAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last = make(map[string]time.Time)
}
