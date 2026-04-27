// Package probe provides per-host probe interval management,
// allowing different scan frequencies to be assigned to individual hosts.
package probe

import (
	"errors"
	"sync"
	"time"
)

const defaultInterval = 60 * time.Second

// Manager stores per-host probe intervals and falls back to a
// configurable default when no specific interval has been set.
type Manager struct {
	mu           sync.RWMutex
	intervals    map[string]time.Duration
	defaultValue time.Duration
}

// New returns a Manager with the given default interval.
// If d is zero, defaultInterval (60 s) is used.
func New(d time.Duration) *Manager {
	if d <= 0 {
		d = defaultInterval
	}
	return &Manager{
		intervals:    make(map[string]time.Duration),
		defaultValue: d,
	}
}

// Set assigns an explicit probe interval for host.
// Returns an error if the interval is not positive.
func (m *Manager) Set(host string, d time.Duration) error {
	if d <= 0 {
		return errors.New("probe: interval must be positive")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.intervals[host] = d
	return nil
}

// Get returns the probe interval for host.
// If no interval has been set the default is returned.
func (m *Manager) Get(host string) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if d, ok := m.intervals[host]; ok {
		return d
	}
	return m.defaultValue
}

// Delete removes the explicit interval for host so that the
// default is used again.
func (m *Manager) Delete(host string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.intervals, host)
}

// All returns a snapshot of every host that has an explicit interval.
func (m *Manager) All() map[string]time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[string]time.Duration, len(m.intervals))
	for k, v := range m.intervals {
		out[k] = v
	}
	return out
}

// Default returns the fallback interval used when no host-specific
// value has been set.
func (m *Manager) Default() time.Duration {
	return m.defaultValue
}
