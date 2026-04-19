package debounce

import (
	"fmt"
	"time"

	"github.com/example/portwatch/internal/snapshot"
)

// Guard wraps a diff-handling func and debounces calls per host+port key.
type Guard struct {
	debouncer *Debouncer
}

// NewGuard returns a Guard using the given quiet window.
func NewGuard(window time.Duration) *Guard {
	return &Guard{debouncer: New(window)}
}

// Filter returns only the diff entries not suppressed by the debounce window.
// It mutates a copy of the diff, leaving the original unchanged.
func (g *Guard) Filter(host string, d snapshot.Diff) snapshot.Diff {
	out := snapshot.Diff{}
	for _, p := range d.Opened {
		key := fmt.Sprintf("%s:opened:%d", host, p)
		if g.debouncer.Allow(key) {
			out.Opened = append(out.Opened, p)
		}
	}
	for _, p := range d.Closed {
		key := fmt.Sprintf("%s:closed:%d", host, p)
		if g.debouncer.Allow(key) {
			out.Closed = append(out.Closed, p)
		}
	}
	return out
}
