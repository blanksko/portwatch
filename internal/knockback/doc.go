// Package knockback provides an exponential back-off gate for scan targets.
//
// When a host experiences repeated consecutive scan failures it is temporarily
// blocked from further scans to avoid wasting resources and generating noise.
// The block duration doubles on each new failure cycle, capped at a
// configurable maximum.
//
// Basic usage:
//
//	g := knockback.New()
//
//	if !g.Allow(host) {
//		// skip scan — host is in back-off
//		return
//	}
//
//	err := scan(host)
//	if err != nil {
//		g.RecordFailure(host)
//	} else {
//		g.RecordSuccess(host)
//	}
package knockback
