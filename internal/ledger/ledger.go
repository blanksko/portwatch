// Package ledger tracks cumulative open-port counts per host over time,
// providing a simple running tally that can be queried or reset.
package ledger

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Entry holds the tally for a single host.
type Entry struct {
	Host      string
	OpenPorts int
	LastSeen  time.Time
}

// Ledger accumulates open-port counts keyed by host.
type Ledger struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	now     func() time.Time
}

// New returns a ready-to-use Ledger.
func New() *Ledger {
	return &Ledger{
		entries: make(map[string]*Entry),
		now:     time.Now,
	}
}

// Record updates the tally for the host in results, counting only open ports.
func (l *Ledger) Record(results []scanner.Result) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, r := range results {
		if !r.Open {
			continue
		}
		e, ok := l.entries[r.Host]
		if !ok {
			e = &Entry{Host: r.Host}
			l.entries[r.Host] = e
		}
		e.OpenPorts++
		e.LastSeen = l.now()
	}
}

// Get returns the entry for host, and whether it exists.
func (l *Ledger) Get(host string) (Entry, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	e, ok := l.entries[host]
	if !ok {
		return Entry{}, false
	}
	return *e, true
}

// Reset clears the tally for host. If host is empty, all entries are cleared.
func (l *Ledger) Reset(host string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if host == "" {
		l.entries = make(map[string]*Entry)
		return
	}
	delete(l.entries, host)
}

// All returns a snapshot of all entries.
func (l *Ledger) All() []Entry {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make([]Entry, 0, len(l.entries))
	for _, e := range l.entries {
		out = append(out, *e)
	}
	return out
}
