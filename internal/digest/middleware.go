package digest

import (
	"log"

	"github.com/user/portwatch/internal/scanner"
)

// Guard wraps a downstream handler and skips it when the port fingerprint
// for a host has not changed since the previous scan cycle.
type Guard struct {
	digest  *Digest
	logger  *log.Logger
	next    func(host string, results []scanner.Result) error
}

// NewGuard returns a Guard that calls next only when the scan results for
// host differ from the previously recorded fingerprint.
func NewGuard(d *Digest, logger *log.Logger, next func(string, []scanner.Result) error) *Guard {
	if d == nil {
		panic("digest: Guard requires a non-nil Digest")
	}
	if logger == nil {
		panic("digest: Guard requires a non-nil logger")
	}
	if next == nil {
		panic("digest: Guard requires a non-nil next handler")
	}
	return &Guard{digest: d, logger: logger, next: next}
}

// Handle checks the fingerprint for host. If unchanged, it logs a debug
// message and returns nil. Otherwise it forwards the call to the wrapped
// handler.
func (g *Guard) Handle(host string, results []scanner.Result) error {
	if !g.digest.Changed(host, results) {
		g.logger.Printf("digest: no change for %s — skipping downstream handlers", host)
		return nil
	}
	return g.next(host, results)
}
