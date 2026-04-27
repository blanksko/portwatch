// Package grace provides a graceful shutdown coordinator that waits for
// in-flight scan cycles to complete before allowing the process to exit.
package grace

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Coordinator listens for OS termination signals and coordinates a clean
// shutdown, giving active work a bounded window to finish.
type Coordinator struct {
	mu      sync.Mutex
	active  int
	done    chan struct{}
	log     *log.Logger
	timeout time.Duration
}

// New returns a Coordinator that will wait up to timeout for active work to
// finish after a signal is received. logger must not be nil.
func New(timeout time.Duration, logger *log.Logger) *Coordinator {
	if logger == nil {
		panic("grace: nil logger")
	}
	return &Coordinator{
		done:    make(chan struct{}),
		log:     logger,
		timeout: timeout,
	}
}

// Acquire marks the start of a unit of work. It returns false if shutdown has
// already been requested, in which case the caller must not proceed.
func (c *Coordinator) Acquire() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	select {
	case <-c.done:
		return false
	default:
	}
	c.active++
	return true
}

// Release marks the completion of a unit of work previously started with
// Acquire.
func (c *Coordinator) Release() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.active > 0 {
		c.active--
	}
}

// Active returns the number of currently in-flight work units.
func (c *Coordinator) Active() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.active
}

// Wait blocks until a SIGINT or SIGTERM is received, then waits up to the
// configured timeout for all active work to finish. The returned context is
// cancelled when it is safe to exit.
func (c *Coordinator) Wait(parent context.Context) context.Context {
	ctx, cancel := context.WithCancel(parent)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer cancel()
		select {
		case sig := <-sigs:
			c.log.Printf("grace: received %s, shutting down", sig)
		case <-parent.Done():
			c.log.Printf("grace: parent context cancelled")
		}
		close(c.done)
		c.drain()
	}()

	return ctx
}

// drain polls until all active work finishes or the timeout elapses.
func (c *Coordinator) drain() {
	deadline := time.Now().Add(c.timeout)
	for time.Now().Before(deadline) {
		if c.Active() == 0 {
			c.log.Printf("grace: all work finished, exiting cleanly")
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	c.log.Printf("grace: timeout elapsed with %d active worker(s), forcing exit", c.Active())
}
