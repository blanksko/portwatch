// Package envelope wraps scan results with metadata for downstream
// processing — host identity, scan time, tags, and labels are bundled
// together so pipeline stages receive a self-contained unit of work.
package envelope

import (
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Envelope carries a scan result set alongside contextual metadata.
type Envelope struct {
	// Host is the target that was scanned.
	Host string

	// ScannedAt is the wall-clock time the scan completed.
	ScannedAt time.Time

	// Results holds the individual port scan results.
	Results []scanner.Result

	// Tags is an optional set of string labels attached to the host.
	Tags []string

	// Labels is an optional map of key/value metadata.
	Labels map[string]string

	// ScanID is a unique identifier for this scan cycle (e.g. a UUID or
	// incrementing counter supplied by the caller).
	ScanID string
}

// New constructs an Envelope for the given host and results, stamping
// ScannedAt with the current UTC time.
func New(host string, results []scanner.Result) *Envelope {
	return &Envelope{
		Host:      host,
		ScannedAt: time.Now().UTC(),
		Results:   results,
		Labels:    make(map[string]string),
	}
}

// WithTags returns the envelope with the supplied tags appended.
func (e *Envelope) WithTags(tags ...string) *Envelope {
	e.Tags = append(e.Tags, tags...)
	return e
}

// WithLabel attaches a single key/value label to the envelope.
func (e *Envelope) WithLabel(key, value string) *Envelope {
	if e.Labels == nil {
		e.Labels = make(map[string]string)
	}
	e.Labels[key] = value
	return e
}

// WithScanID sets the ScanID field and returns the envelope for chaining.
func (e *Envelope) WithScanID(id string) *Envelope {
	e.ScanID = id
	return e
}

// OpenCount returns the number of results whose Open field is true.
func (e *Envelope) OpenCount() int {
	n := 0
	for _, r := range e.Results {
		if r.Open {
			n++
		}
	}
	return n
}
