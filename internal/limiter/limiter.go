// Package limiter provides a concurrency limiter that caps the number of
// simultaneous port scans running against any single host.
package limiter

import (
	"fmt"
	"sync"
)

const defaultConcurrency = 10

// Limiter controls how many goroutines may scan a given host at once.
type Limiter struct {
	mu      sync.Mutex
	sems    map[string]chan struct{}
	limit   int
}

// New returns a Limiter that allows at most concurrency simultaneous
// acquisitions per host. If concurrency is zero or negative the default
// of 10 is used.
func New(concurrency int) *Limiter {
	if concurrency <= 0 {
		concurrency = defaultConcurrency
	}
	return &Limiter{
		sems:  make(map[string]chan struct{}),
		limit: concurrency,
	}
}

// Acquire blocks until a slot is available for host, then claims it.
// The caller must call Release with the same host when done.
func (l *Limiter) Acquire(host string) {
	l.mu.Lock()
	ch, ok := l.sems[host]
	if !ok {
		ch = make(chan struct{}, l.limit)
		l.sems[host] = ch
	}
	l.mu.Unlock()
	ch <- struct{}{}
}

// Release frees one slot previously claimed by Acquire for host.
// It panics if host has no outstanding acquisitions.
func (l *Limiter) Release(host string) {
	l.mu.Lock()
	ch, ok := l.sems[host]
	l.mu.Unlock()
	if !ok {
		panic(fmt.Sprintf("limiter: Release called for unknown host %q", host))
	}
	<-ch
}

// Active returns the number of currently held slots for host.
func (l *Limiter) Active(host string) int {
	l.mu.Lock()
	ch, ok := l.sems[host]
	l.mu.Unlock()
	if !ok {
		return 0
	}
	return len(ch)
}

// Reset removes all tracking state for host, releasing any buffered slots.
func (l *Limiter) Reset(host string) {
	l.mu.Lock()
	delete(l.sems, host)
	l.mu.Unlock()
}
