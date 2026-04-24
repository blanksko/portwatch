// Package dedup provides result deduplication to suppress identical
// consecutive scan outputs for a given host, reducing downstream noise.
package dedup

import (
	"sync"

	"github.com/yourorg/portwatch/internal/scanner"
)

// Store tracks the last seen set of open ports per host so that callers
// can skip processing when nothing has changed since the previous scan.
type Store struct {
	mu   sync.Mutex
	last map[string]string // host -> comma-joined sorted port fingerprint
}

// New returns an initialised Store.
func New() *Store {
	return &Store{last: make(map[string]string)}
}

// Changed reports whether the results differ from the last recorded state
// for the same host. It always returns true on the first call for a host.
// If results is empty the host key is left unchanged.
func (s *Store) Changed(host string, results []scanner.Result) bool {
	if len(results) == 0 {
		return false
	}
	fp := fingerprint(results)
	s.mu.Lock()
	defer s.mu.Unlock()
	prev, seen := s.last[host]
	if !seen || prev != fp {
		s.last[host] = fp
		return true
	}
	return false
}

// Reset removes the cached fingerprint for host so the next call to
// Changed will always return true.
func (s *Store) Reset(host string) {
	s.mu.Lock()
	delete(s.last, host)
	s.mu.Unlock()
}

// fingerprint returns a stable string representation of the open ports
// present in results.
func fingerprint(results []scanner.Result) string {
	ports := make([]byte, 0, len(results)*6)
	for i, r := range results {
		if !r.Open {
			continue
		}
		if i > 0 {
			ports = append(ports, ',')
		}
		ports = appendPort(ports, r.Port)
	}
	return string(ports)
}

func appendPort(b []byte, port int) []byte {
	if port == 0 {
		return append(b, '0')
	}
	var buf [6]byte
	i := len(buf)
	for port > 0 {
		i--
		buf[i] = byte('0' + port%10)
		port /= 10
	}
	return append(b, buf[i:]...)
}
