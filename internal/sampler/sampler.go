// Package sampler provides port-scan result sampling to reduce noise
// by collecting multiple scan results and returning a representative subset.
package sampler

import (
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Sampler collects scan results over a window and returns a deduplicated,
// stable sample of observed open ports per host.
type Sampler struct {
	mu      sync.Mutex
	buckets map[string]map[int]int // host -> port -> seen count
	thresh  int
}

// New creates a Sampler that promotes a port to the sample once it has been
// observed at least thresh times across successive scans.
func New(thresh int) *Sampler {
	if thresh < 1 {
		thresh = 1
	}
	return &Sampler{
		buckets: make(map[string]map[int]int),
		thresh:  thresh,
	}
}

// Record adds a batch of scan results to the internal counters.
func (s *Sampler) Record(results []scanner.Result) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, r := range results {
		if !r.Open {
			continue
		}
		if _, ok := s.buckets[r.Host]; !ok {
			s.buckets[r.Host] = make(map[int]int)
		}
		s.buckets[r.Host][r.Port]++
	}
}

// Sample returns only the results whose ports have been seen at least thresh
// times, providing a stable view of consistently open ports.
func (s *Sampler) Sample(host string) []scanner.Result {
	s.mu.Lock()
	defer s.mu.Unlock()

	ports, ok := s.buckets[host]
	if !ok {
		return nil
	}

	var out []scanner.Result
	for port, count := range ports {
		if count >= s.thresh {
			out = append(out, scanner.Result{Host: host, Port: port, Open: true})
		}
	}
	return out
}

// Reset clears all accumulated data for a given host.
func (s *Sampler) Reset(host string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.buckets, host)
}
