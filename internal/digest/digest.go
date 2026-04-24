// Package digest computes and compares fingerprints of scan results,
// allowing portwatch to detect when a host's port profile has changed
// without storing a full snapshot.
package digest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Digest holds the last-known fingerprint for each host.
type Digest struct {
	mu     sync.RWMutex
	store  map[string]string
}

// New returns an initialised Digest store.
func New() *Digest {
	return &Digest{store: make(map[string]string)}
}

// Compute returns a deterministic hex fingerprint for the given scan results.
// The fingerprint is derived from the sorted list of open port numbers.
func Compute(results []scanner.Result) string {
	ports := make([]int, 0, len(results))
	for _, r := range results {
		if r.Open {
			ports = append(ports, r.Port)
		}
	}
	sort.Ints(ports)

	h := sha256.New()
	for _, p := range ports {
		fmt.Fprintf(h, "%d\n", p)
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Changed reports whether the fingerprint for host differs from the last
// recorded value. It also stores the new fingerprint.
func (d *Digest) Changed(host string, results []scanner.Result) bool {
	next := Compute(results)

	d.mu.Lock()
	defer d.mu.Unlock()

	prev, ok := d.store[host]
	d.store[host] = next
	return !ok || prev != next
}

// Get returns the stored fingerprint for host, or an empty string if none
// has been recorded yet.
func (d *Digest) Get(host string) string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.store[host]
}

// Reset removes the stored fingerprint for host, causing the next call to
// Changed to always return true.
func (d *Digest) Reset(host string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.store, host)
}
