package notifier

import (
	"fmt"
	"io"
	"strings"

	"github.com/user/portwatch/internal/snapshot"
)

// Notifier formats and writes diff summaries to an output destination.
type Notifier struct {
	w io.Writer
}

// New creates a Notifier that writes to w.
// It returns an error if w is nil.
func New(w io.Writer) (*Notifier, error) {
	if w == nil {
		return nil, fmt.Errorf("notifier: writer must not be nil")
	}
	return &Notifier{w: w}, nil
}

// Notify writes a human-readable summary of d to the underlying writer.
// It returns nil when there are no changes.
func (n *Notifier) Notify(host string, d snapshot.DiffResult) error {
	if len(d.Opened) == 0 && len(d.Closed) == 0 {
		return nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("[portwatch] changes detected on %s\n", host))

	if len(d.Opened) > 0 {
		sb.WriteString(fmt.Sprintf("  opened ports (%d):\n", len(d.Opened)))
		for _, p := range d.Opened {
			sb.WriteString(fmt.Sprintf("    + %d\n", p))
		}
	}

	if len(d.Closed) > 0 {
		sb.WriteString(fmt.Sprintf("  closed ports (%d):\n", len(d.Closed)))
		for _, p := range d.Closed {
			sb.WriteString(fmt.Sprintf("    - %d\n", p))
		}
	}

	_, err := fmt.Fprint(n.w, sb.String())
	return err
}
