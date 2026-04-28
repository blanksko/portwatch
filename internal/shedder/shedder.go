// Package shedder implements load shedding for port scan jobs.
// When the number of queued or active scans exceeds a configured
// threshold, new scan requests are dropped to protect the host.
package shedder

import (
	"errors"
	"sync/atomic"
)

// ErrShed is returned when a scan request is dropped by the shedder.
var ErrShed = errors.New("shedder: load too high, request dropped")

// Shedder tracks active work and rejects new requests above a limit.
type Shedder struct {
	limit   int64
	active  atomic.Int64
}

// New returns a Shedder that allows at most limit concurrent scans.
// Panics if limit is less than 1.
func New(limit int) *Shedder {
	if limit < 1 {
		panic("shedder: limit must be at least 1")
	}
	return &Shedder{limit: int64(limit)}
}

// Acquire attempts to reserve a slot for a new scan.
// Returns ErrShed if the active count is at or above the limit.
func (s *Shedder) Acquire() error {
	if s.active.Load() >= s.limit {
		return ErrShed
	}
	s.active.Add(1)
	return nil
}

// Release decrements the active counter. It must be called exactly
// once for every successful Acquire.
func (s *Shedder) Release() {
	if v := s.active.Add(-1); v < 0 {
		s.active.Store(0)
	}
}

// Active returns the current number of in-flight scans.
func (s *Shedder) Active() int {
	return int(s.active.Load())
}

// Limit returns the configured concurrency ceiling.
func (s *Shedder) Limit() int {
	return int(s.limit)
}
