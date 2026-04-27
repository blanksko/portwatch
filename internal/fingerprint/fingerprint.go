// Package fingerprint produces a stable string identity for a host's
// current set of open ports. Two scans that yield identical open ports
// on the same host will always produce the same fingerprint, regardless
// of the order in which results were returned by the scanner.
package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"

	"github.com/user/portwatch/internal/scanner"
)

// Fingerprint is a hex-encoded SHA-256 digest that uniquely identifies
// the set of open ports observed for a host at a point in time.
type Fingerprint string

// Compute derives a deterministic Fingerprint from a slice of scan results.
// Only open ports are included; closed ports are ignored.
// Results belonging to different hosts are each fingerprinted independently
// and the host name is mixed into the digest so that identical port sets on
// different hosts produce different fingerprints.
func Compute(results []scanner.Result) Fingerprint {
	if len(results) == 0 {
		return Fingerprint(empty())
	}

	// Collect and sort open-port tokens for determinism.
	tokens := make([]string, 0, len(results))
	for _, r := range results {
		if !r.Open {
			continue
		}
		tokens = append(tokens, fmt.Sprintf("%s:%d", r.Host, r.Port))
	}

	if len(tokens) == 0 {
		return Fingerprint(empty())
	}

	sort.Strings(tokens)

	h := sha256.New()
	for _, t := range tokens {
		_, _ = fmt.Fprintln(h, t)
	}
	return Fingerprint(hex.EncodeToString(h.Sum(nil)))
}

// Equal reports whether two fingerprints are identical.
func Equal(a, b Fingerprint) bool { return a == b }

// String returns the hex string representation of the fingerprint.
func (f Fingerprint) String() string { return string(f) }

// empty returns the SHA-256 digest of an empty input.
func empty() string {
	h := sha256.New()
	return hex.EncodeToString(h.Sum(nil))
}
