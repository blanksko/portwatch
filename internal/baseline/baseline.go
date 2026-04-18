// Package baseline manages the expected (approved) port state for a host.
// It allows operators to "accept" the current snapshot as the known-good baseline
// and later compare live scans against it.
package baseline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Baseline holds an approved set of scan results for a host.
type Baseline struct {
	Host      string           `json:"host"`
	Ports     []scanner.Result `json:"ports"`
	CreatedAt time.Time        `json:"created_at"`
}

// Manager persists and loads baselines from a directory.
type Manager struct {
	dir string
}

// New returns a Manager that stores baselines under dir.
func New(dir string) (*Manager, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("baseline: create dir: %w", err)
	}
	return &Manager{dir: dir}, nil
}

// Save writes the baseline for the given host to disk.
func (m *Manager) Save(host string, results []scanner.Result) error {
	b := Baseline{Host: host, Ports: results, CreatedAt: time.Now().UTC()}
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("baseline: marshal: %w", err)
	}
	return os.WriteFile(m.filePath(host), data, 0o644)
}

// Load retrieves the stored baseline for host. Returns nil, nil if none exists.
func (m *Manager) Load(host string) (*Baseline, error) {
	data, err := os.ReadFile(m.filePath(host))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("baseline: read: %w", err)
	}
	var b Baseline
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("baseline: unmarshal: %w", err)
	}
	return &b, nil
}

// Delete removes the baseline file for host.
func (m *Manager) Delete(host string) error {
	err := os.Remove(m.filePath(host))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (m *Manager) filePath(host string) string {
	safe := sanitize(host)
	return filepath.Join(m.dir, safe+".baseline.json")
}

func sanitize(host string) string {
	out := make([]byte, len(host))
	for i := range host {
		if host[i] == ':' || host[i] == '/' || host[i] == '\\' {
			out[i] = '_'
		} else {
			out[i] = host[i]
		}
	}
	return string(out)
}
