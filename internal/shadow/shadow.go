// Package shadow provides a secondary passive scan store that records
// scan results without triggering alerts, useful for baselining and
// comparing against active scan output.
package shadow

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Record holds a shadow scan entry for a single host.
type Record struct {
	Host      string           `json:"host"`
	Results   []scanner.Result `json:"results"`
	RecordedAt time.Time       `json:"recorded_at"`
}

// Store persists shadow scan records to disk.
type Store struct {
	mu  sync.RWMutex
	dir string
}

// New creates a new Store rooted at dir, creating the directory if needed.
func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("shadow: create dir: %w", err)
	}
	return &Store{dir: dir}, nil
}

// Record saves a shadow snapshot for the given host.
func (s *Store) Record(host string, results []scanner.Result) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	rec := Record{
		Host:       host,
		Results:    results,
		RecordedAt: time.Now().UTC(),
	}
	data, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("shadow: marshal: %w", err)
	}
	return os.WriteFile(s.path(host), data, 0o644)
}

// Load retrieves the shadow record for the given host.
// Returns os.ErrNotExist if no record has been saved yet.
func (s *Store) Load(host string) (Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := os.ReadFile(s.path(host))
	if err != nil {
		return Record{}, err
	}
	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return Record{}, fmt.Errorf("shadow: unmarshal: %w", err)
	}
	return rec, nil
}

// Delete removes the shadow record for host, if present.
func (s *Store) Delete(host string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := os.Remove(s.path(host))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *Store) path(host string) string {
	safe := sanitize(host)
	return filepath.Join(s.dir, safe+".json")
}

func sanitize(host string) string {
	out := make([]byte, len(host))
	for i := range host {
		switch host[i] {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
			out[i] = '_'
		default:
			out[i] = host[i]
		}
	}
	return string(out)
}
