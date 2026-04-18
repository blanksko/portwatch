// Package audit provides structured audit logging of port scan events.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Host      string    `json:"host"`
	Opened    []int     `json:"opened,omitempty"`
	Closed    []int     `json:"closed,omitempty"`
	Message   string    `json:"message"`
}

// Logger writes audit entries as newline-delimited JSON.
type Logger struct {
	w io.Writer
}

// New returns a new audit Logger writing to w.
// It panics if w is nil.
func New(w io.Writer) *Logger {
	if w == nil {
		panic("audit: writer must not be nil")
	}
	return &Logger{w: w}
}

// Record writes an audit entry for the given host and diff.
func (l *Logger) Record(host string, d snapshot.Diff) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Host:      host,
		Opened:    d.Opened,
		Closed:    d.Closed,
		Message:   buildMessage(d),
	}
	b, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", b)
	return err
}

func buildMessage(d snapshot.Diff) string {
	switch {
	case len(d.Opened) > 0 && len(d.Closed) > 0:
		return "ports opened and closed"
	case len(d.Opened) > 0:
		return "new ports detected"
	case len(d.Closed) > 0:
		return "ports closed"
	default:
		return "no changes"
	}
}
