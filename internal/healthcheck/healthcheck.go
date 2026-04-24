// Package healthcheck provides per-host liveness probing to distinguish
// unresponsive hosts from hosts with no open ports.
package healthcheck

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const defaultTimeout = 2 * time.Second

// Status represents the liveness state of a host.
type Status int

const (
	StatusUnknown Status = iota
	StatusAlive
	StatusUnreachable
)

func (s Status) String() string {
	switch s {
	case StatusAlive:
		return "alive"
	case StatusUnreachable:
		return "unreachable"
	default:
		return "unknown"
	}
}

// Checker probes hosts for basic TCP reachability.
type Checker struct {
	timeout time.Duration
	logger  *log.Logger
	mu      sync.RWMutex
	cache   map[string]Status
}

// New returns a Checker with the given timeout. If timeout is zero the default
// is used. logger must not be nil.
func New(timeout time.Duration, logger *log.Logger) *Checker {
	if logger == nil {
		panic("healthcheck: logger must not be nil")
	}
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return &Checker{
		timeout: timeout,
		logger:  logger,
		cache:   make(map[string]Status),
	}
}

// Probe checks whether host is reachable by attempting a TCP connection to
// port 80 or 443. The result is cached and returned.
func (c *Checker) Probe(ctx context.Context, host string) Status {
	for _, port := range []string{"80", "443"} {
		addr := net.JoinHostPort(host, port)
		dialer := &net.Dialer{Timeout: c.timeout}
		conn, err := dialer.DialContext(ctx, "tcp", addr)
		if err == nil {
			conn.Close()
			c.setCache(host, StatusAlive)
			return StatusAlive
		}
	}
	c.logger.Printf("healthcheck: host %s unreachable", host)
	c.setCache(host, StatusUnreachable)
	return StatusUnreachable
}

// Last returns the most recently cached status for host.
func (c *Checker) Last(host string) Status {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if s, ok := c.cache[host]; ok {
		return s
	}
	return StatusUnknown
}

// Reset clears the cached status for host.
func (c *Checker) Reset(host string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, fmt.Sprintf("%s", host))
}

func (c *Checker) setCache(host string, s Status) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[host] = s
}
