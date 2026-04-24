package sampler

import (
	"context"
	"log"

	"github.com/user/portwatch/internal/scanner"
)

// ScanFunc is the signature of a downstream scan handler.
type ScanFunc func(ctx context.Context, host string, results []scanner.Result) error

// Guard wraps a ScanFunc so that results are recorded in the Sampler before
// being forwarded. Only ports that have met the frequency threshold are passed
// downstream, silently dropping transient observations.
type Guard struct {
	s      *Sampler
	next   ScanFunc
	logger *log.Logger
}

// NewGuard creates a Guard that filters results through s before calling next.
func NewGuard(s *Sampler, next ScanFunc, logger *log.Logger) *Guard {
	if s == nil {
		panic("sampler: Guard requires a non-nil Sampler")
	}
	if next == nil {
		panic("sampler: Guard requires a non-nil next handler")
	}
	if logger == nil {
		panic("sampler: Guard requires a non-nil logger")
	}
	return &Guard{s: s, next: next, logger: logger}
}

// Handle records results and forwards only the stable sample downstream.
func (g *Guard) Handle(ctx context.Context, host string, results []scanner.Result) error {
	g.s.Record(results)

	stable := g.s.Sample(host)
	if len(stable) == 0 {
		g.logger.Printf("sampler: no stable ports yet for %s, skipping downstream", host)
		return nil
	}

	return g.next(ctx, host, stable)
}
