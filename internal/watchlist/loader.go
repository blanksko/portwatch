package watchlist

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// FromReader populates a Watchlist from a plain-text reader. Each non-empty,
// non-comment line is expected to have the form:
//
//	<host> [alias]
//
// Lines beginning with '#' are treated as comments and skipped.
func FromReader(r io.Reader) (*Watchlist, error) {
	w := New()
	scanner := bufio.NewScanner(r)
	line := 0
	for scanner.Scan() {
		line++
		raw := strings.TrimSpace(scanner.Text())
		if raw == "" || strings.HasPrefix(raw, "#") {
			continue
		}
		parts := strings.Fields(raw)
		host := parts[0]
		alias := ""
		if len(parts) >= 2 {
			alias = parts[1]
		}
		if err := w.Add(host, alias); err != nil {
			return nil, fmt.Errorf("watchlist: line %d: %w", line, err)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("watchlist: read error: %w", err)
	}
	return w, nil
}

// FromStrings is a convenience wrapper that builds a Watchlist from a slice
// of host strings (no aliases).
func FromStrings(hosts []string) (*Watchlist, error) {
	w := New()
	for _, h := range hosts {
		if err := w.Add(h, ""); err != nil {
			return nil, err
		}
	}
	return w, nil
}
