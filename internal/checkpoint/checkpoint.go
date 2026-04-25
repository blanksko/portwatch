// Package checkpoint persists the last-seen scan time for each host so that
// portwatch can resume monitoring without re-alerting on already-known state
// after a restart.
package checkpoint

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Store persists and retrieves per-host checkpoint timestamps.
type Store struct {
	mu   sync.RWMutex
	dir  string
	data map[string]time.Time
}

// New creates a Store rooted at dir, loading any previously saved state.
func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	s := &Store{dir: dir, data: make(map[string]time.Time)}
	if err := s.load(); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	return s, nil
}

// Set records the checkpoint time for host.
func (s *Store) Set(host string, t time.Time) error {
	s.mu.Lock()
	s.data[host] = t
	s.mu.Unlock()
	return s.flush()
}

// Get returns the last checkpoint time for host and whether it was found.
func (s *Store) Get(host string) (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.data[host]
	return t, ok
}

// Delete removes the checkpoint entry for host and persists the change.
func (s *Store) Delete(host string) error {
	s.mu.Lock()
	delete(s.data, host)
	s.mu.Unlock()
	return s.flush()
}

// All returns a snapshot of all stored checkpoints.
func (s *Store) All() map[string]time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]time.Time, len(s.data))
	for k, v := range s.data {
		out[k] = v
	}
	return out
}

func (s *Store) path() string {
	return filepath.Join(s.dir, "checkpoints.json")
}

func (s *Store) load() error {
	b, err := os.ReadFile(s.path())
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return json.Unmarshal(b, &s.data)
}

func (s *Store) flush() error {
	s.mu.RLock()
	b, err := json.MarshalIndent(s.data, "", "  ")
	s.mu.RUnlock()
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(), b, 0o644)
}
