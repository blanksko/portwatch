// Package snapshot provides functionality for saving and loading port scan
// snapshots to disk, enabling change detection between runs.
package snapshot

import (
	"encoding/json"
	"os"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot represents a saved state of a port scan.
type Snapshot struct {
	Host      string                `json:"host"`
	Timestamp time.Time             `json:"timestamp"`
	Results   []scanner.ScanResult  `json:"results"`
}

// New creates a new Snapshot from scan results.
func New(host string, results []scanner.ScanResult) *Snapshot {
	return &Snapshot{
		Host:      host,
		Timestamp: time.Now(),
		Results:   results,
	}
}

// Save writes the snapshot to the given file path as JSON.
func Save(path string, snap *Snapshot) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}

// Load reads a snapshot from the given file path.
func Load(path string) (*Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var snap Snapshot
	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return nil, err
	}
	return &snap, nil
}
