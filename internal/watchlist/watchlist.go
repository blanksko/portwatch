package watchlist

import (
	"fmt"
	"sync"
)

// Entry represents a single watched host with an optional alias.
type Entry struct {
	Host  string
	Alias string
}

// Watchlist maintains an ordered, deduplicated set of hosts to monitor.
type Watchlist struct {
	mu      sync.RWMutex
	entries []Entry
	index   map[string]int // host -> slice position
}

// New returns an empty Watchlist.
func New() *Watchlist {
	return &Watchlist{
		index: make(map[string]int),
	}
}

// Add inserts a host into the watchlist. If the host already exists the call
// is a no-op. An optional alias may be supplied for display purposes.
func (w *Watchlist) Add(host, alias string) error {
	if host == "" {
		return fmt.Errorf("watchlist: host must not be empty")
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if _, ok := w.index[host]; ok {
		return nil
	}
	w.index[host] = len(w.entries)
	w.entries = append(w.entries, Entry{Host: host, Alias: alias})
	return nil
}

// Remove deletes a host from the watchlist. It returns false when the host
// was not present.
func (w *Watchlist) Remove(host string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	idx, ok := w.index[host]
	if !ok {
		return false
	}
	w.entries = append(w.entries[:idx], w.entries[idx+1:]...)
	delete(w.index, host)
	// Rebuild index for shifted entries.
	for i := idx; i < len(w.entries); i++ {
		w.index[w.entries[i].Host] = i
	}
	return true
}

// Has reports whether host is currently in the watchlist.
func (w *Watchlist) Has(host string) bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	_, ok := w.index[host]
	return ok
}

// All returns a snapshot of the current entries in insertion order.
func (w *Watchlist) All() []Entry {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]Entry, len(w.entries))
	copy(out, w.entries)
	return out
}

// Len returns the number of watched hosts.
func (w *Watchlist) Len() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.entries)
}
