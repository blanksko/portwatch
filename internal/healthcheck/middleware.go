package healthcheck

import (
	"context"
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// Guard is a scan middleware that skips downstream processing when a host is
// determined to be unreachable, preventing false-positive "all ports closed"
// alerts during outages.
type Guard struct {
	checker *Checker
	next    func(ctx context.Context, host string, results []scanner.Result) error
}

// NewGuard wraps next with a liveness check. If the host is unreachable the
// guard returns nil without calling next. checker must not be nil.
func NewGuard(
	checker *Checker,
	next func(ctx context.Context, host string, results []scanner.Result) error,
) *Guard {
	if checker == nil {
		panic("healthcheck: Guard checker must not be nil")
	}
	if next == nil {
		panic("healthcheck: Guard next must not be nil")
	}
	return &Guard{checker: checker, next: next}
}

// Run probes host and, if alive, delegates to the wrapped handler.
func (g *Guard) Run(ctx context.Context, host string, results []scanner.Result) error {
	status := g.checker.Probe(ctx, host)
	if status == StatusUnreachable {
		return fmt.Errorf("healthcheck: %s is unreachable, skipping", host)
	}
	return g.next(ctx, host, results)
}
