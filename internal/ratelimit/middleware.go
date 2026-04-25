package ratelimit

import (
	"fmt"
	"log"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Guard wraps a scan job and skips execution for a host if the rate limit
// has not yet elapsed since the last allowed call.
type Guard struct {
	limiter *RateLimit
	logger  *log.Logger
	next    func(host string) ([]scanner.Result, error)
}

// NewGuard returns a Guard that enforces r on every host before delegating
// to next. Panics if r or logger is nil.
func NewGuard(r *RateLimit, logger *log.Logger, next func(host string) ([]scanner.Result, error)) *Guard {
	if r == nil {
		panic("ratelimit: Guard requires a non-nil RateLimit")
	}
	if logger == nil {
		panic("ratelimit: Guard requires a non-nil logger")
	}
	return &Guard{limiter: r, logger: logger, next: next}
}

// Run checks the rate limit for host. If the limit has not elapsed it returns
// a nil slice and no error. Otherwise it delegates to the wrapped function.
func (g *Guard) Run(host string) ([]scanner.Result, error) {
	if !g.limiter.Allow(host) {
		g.logger.Printf("ratelimit: skipping %s — rate limit active", host)
		return nil, nil
	}
	return g.next(host)
}

// Reset clears the rate-limit state for host, allowing the next call
// through immediately regardless of the configured interval.
func (g *Guard) Reset(host string) {
	g.limiter.Reset(host)
}

// String returns a human-readable description of the guard's interval.
func (g *Guard) String() string {
	return fmt.Sprintf("RateLimitGuard(interval=%s)", g.limiter.interval.Round(time.Millisecond))
}
