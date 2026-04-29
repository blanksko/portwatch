// Package fence provides a port-range boundary guard that restricts
// scan results to a declared set of allowed port ranges.
package fence

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Range describes an inclusive [Low, High] port boundary.
type Range struct {
	Low  int
	High int
}

// Fence holds a collection of allowed port ranges and filters scan results
// against them.
type Fence struct {
	mu     sync.RWMutex
	ranges []Range
}

// New returns a Fence initialised with the supplied ranges.
// Ranges with Low > High are rejected.
func New(ranges []Range) (*Fence, error) {
	for _, r := range ranges {
		if r.Low > r.High {
			return nil, fmt.Errorf("fence: invalid range %d-%d: low exceeds high", r.Low, r.High)
		}
		if r.Low < 1 || r.High > 65535 {
			return nil, fmt.Errorf("fence: range %d-%d out of valid port bounds (1-65535)", r.Low, r.High)
		}
	}
	return &Fence{ranges: append([]Range(nil), ranges...)}, nil
}

// Add appends a new range to the fence at runtime.
func (f *Fence) Add(r Range) error {
	if r.Low > r.High {
		return fmt.Errorf("fence: invalid range %d-%d", r.Low, r.High)
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.ranges = append(f.ranges, r)
	return nil
}

// Within reports whether port p falls inside any declared range.
func (f *Fence) Within(p int) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	for _, r := range f.ranges {
		if p >= r.Low && p <= r.High {
			return true
		}
	}
	return false
}

// Apply returns only those results whose port falls within the fence.
// Results with a closed port are always dropped.
func (f *Fence) Apply(results []scanner.Result) []scanner.Result {
	out := make([]scanner.Result, 0, len(results))
	for _, res := range results {
		if res.Open && f.Within(res.Port) {
			out = append(out, res)
		}
	}
	return out
}

// Ranges returns a snapshot of the currently configured ranges.
func (f *Fence) Ranges() []Range {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return append([]Range(nil), f.ranges...)
}
