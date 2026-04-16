// Package reporter formats and writes scan reports to an output destination.
package reporter

import (
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Reporter writes human-readable scan reports.
type Reporter struct {
	w io.Writer
}

// New creates a new Reporter that writes to w.
func New(w io.Writer) (*Reporter, error) {
	if w == nil {
		return nil, fmt.Errorf("reporter: writer must not be nil")
	}
	return &Reporter{w: w}, nil
}

// Report writes a formatted summary of a diff to the reporter's writer.
func (r *Reporter) Report(d snapshot.Diff) error {
	timestamp := time.Now().Format(time.RFC3339)

	fmt.Fprintf(r.w, "[%s] Port scan report\n", timestamp)

	if len(d.Opened) == 0 && len(d.Closed) == 0 {
		fmt.Fprintln(r.w, "  No changes detected.")
		return nil
	}

	if len(d.Opened) > 0 {
		fmt.Fprintf(r.w, "  Opened ports (%d):\n", len(d.Opened))
		for _, p := range d.Opened {
			fmt.Fprintf(r.w, "    + %d/%s (%s)\n", p.Port, p.Protocol, p.Service)
		}
	}

	if len(d.Closed) > 0 {
		fmt.Fprintf(r.w, "  Closed ports (%d):\n", len(d.Closed))
		for _, p := range d.Closed {
			fmt.Fprintf(r.w, "    - %d/%s (%s)\n", p.Port, p.Protocol, p.Service)
		}
	}

	return nil
}
