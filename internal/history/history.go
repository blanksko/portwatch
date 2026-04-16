// Package history provides persistent storage of scan snapshots over time,
// allowing portwatch to maintain a log of historical port state changes.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry represents a single historical scan record.
type Entry struct {
	Timestamp time.Time        `json:"timestamp"`
	Host      string           `json:"host"`
	Results   []scanner.Result `json:"results"`
}

// History manages a log of scan entries on disk.
type History struct {
	dir string
}

// New returns a History that stores entries under dir.
func New(dir string) (*History, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("history: create dir: %w", err)
	}
	return &History{dir: dir}, nil
}

// Record writes a new entry to the history log for the given host.
func (h *History) Record(host string, results []scanner.Result) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Host:      host,
		Results:   results,
	}
	filename := fmt.Sprintf("%s_%s.json", sanitize(host), entry.Timestamp.Format("20060102T150405Z"))
	path := filepath.Join(h.dir, filename)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("history: create file: %w", err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(entry); err != nil {
		return fmt.Errorf("history: encode entry: %w", err)
	}
	return nil
}

// sanitize replaces characters unsafe for filenames with underscores.
func sanitize(s string) string {
	out := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '/' || c == ':' || c == '\\' {
			out[i] = '_'
		} else {
			out[i] = c
		}
	}
	return string(out)
}
