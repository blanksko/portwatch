// Package timeout provides per-host scan timeout management.
package timeout

import (
	"sync"
	"time"
)

// Manager holds timeout durations keyed by host.
type Manager struct {
	mu      sync.RWMutex
	default_ time.Duration
	overrides map[string]time.Duration
}

// New returns a Manager with the given default timeout.
func New(defaultTimeout time.Duration) *Manager {
	if defaultTimeout <= 0 {
		defaultTimeout = 5 * time.Second
	}
	return &Manager{
		default_:  defaultTimeout,
		overrides: make(map[string]time.Duration),
	}
}

// Set registers a custom timeout for a specific host.
func (m *Manager) Set(host string, d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.overrides[host] = d
}

// Get returns the timeout for the given host, falling back to the default.
func (m *Manager) Get(host string) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if d, ok := m.overrides[host]; ok {
		return d
	}
	return m.default_
}

// Delete removes a host override, reverting it to the default.
func (m *Manager) Delete(host string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.overrides, host)
}

// Default returns the configured default timeout.
func (m *Manager) Default() time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.default_
}
