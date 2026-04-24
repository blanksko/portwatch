// Package rollup aggregates multiple scan diffs over a time window
// and emits a single combined diff, reducing alert noise during
// periods of rapid port churn.
package rollup

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Rollup accumulates diffs within a fixed window and returns a
// merged view when the window closes.
type Rollup struct {
	mu       sync.Mutex
	window   time.Duration
	opened   map[int]struct{}
	closed   map[int]struct{}
	firstAt  time.Time
	now      func() time.Time
}

// New returns a Rollup that merges diffs arriving within window.
func New(window time.Duration) *Rollup {
	return &Rollup{
		window: window,
		opened: make(map[int]struct{}),
		closed: make(map[int]struct{}),
		now:    time.Now,
	}
}

// Add incorporates a diff into the current accumulation window.
// It returns (merged, true) when the window has elapsed, otherwise
// returns (zero, false).
func (r *Rollup) Add(d snapshot.Diff) (snapshot.Diff, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := r.now()

	if r.firstAt.IsZero() {
		r.firstAt = now
	}

	for _, p := range d.Opened {
		// A port opened then closed in same window cancels out.
		delete(r.closed, p)
		r.opened[p] = struct{}{}
	}
	for _, p := range d.Closed {
		delete(r.opened, p)
		r.closed[p] = struct{}{}
	}

	if now.Sub(r.firstAt) < r.window {
		return snapshot.Diff{}, false
	}

	merged := r.flush()
	return merged, true
}

// Flush forces emission of whatever has been accumulated so far,
// resetting internal state regardless of the window.
func (r *Rollup) Flush() snapshot.Diff {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.flush()
}

func (r *Rollup) flush() snapshot.Diff {
	var d snapshot.Diff
	for p := range r.opened {
		d.Opened = append(d.Opened, p)
	}
	for p := range r.closed {
		d.Closed = append(d.Closed, p)
	}
	r.opened = make(map[int]struct{})
	r.closed = make(map[int]struct{})
	r.firstAt = time.Time{}
	return d
}
