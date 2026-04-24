package dedup

import (
	"context"
	"log/slog"

	"github.com/yourorg/portwatch/internal/scanner"
)

// NextFunc is the signature of the downstream handler that receives results.
type NextFunc func(ctx context.Context, host string, results []scanner.Result) error

// Guard wraps a NextFunc and skips it when results are identical to the
// previous scan for the same host.
type Guard struct {
	store  *Store
	next   NextFunc
	logger *slog.Logger
}

// NewGuard returns a Guard that gates next behind the provided Store.
// Panics if store, next, or logger is nil.
func NewGuard(store *Store, next NextFunc, logger *slog.Logger) *Guard {
	if store == nil {
		panic("dedup: store must not be nil")
	}
	if next == nil {
		panic("dedup: next must not be nil")
	}
	if logger == nil {
		panic("dedup: logger must not be nil")
	}
	return &Guard{store: store, next: next, logger: logger}
}

// Handle checks whether results have changed for host. If they have it
// forwards the call to the wrapped NextFunc; otherwise it logs a debug
// message and returns nil.
func (g *Guard) Handle(ctx context.Context, host string, results []scanner.Result) error {
	if !g.store.Changed(host, results) {
		g.logger.Debug("dedup: skipping unchanged results", "host", host)
		return nil
	}
	return g.next(ctx, host, results)
}
