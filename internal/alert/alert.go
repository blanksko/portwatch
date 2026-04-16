// Package alert provides functionality for reporting port change alerts
// to various outputs such as stdout or log files.
package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Alert represents a single port change notification.
type Alert struct {
	Timestamp time.Time
	Level     Level
	Message   string
}

// Notifier sends alerts to a given writer.
type Notifier struct {
	w io.Writer
}

// New creates a Notifier that writes to w. If w is nil, os.Stdout is used.
func New(w io.Writer) *Notifier {
	if w == nil {
		w = os.Stdout
	}
	return &Notifier{w: w}
}

// Notify formats and writes alerts based on a snapshot diff result.
func (n *Notifier) Notify(d snapshot.DiffResult) []Alert {
	var alerts []Alert
	now := time.Now()

	for _, p := range d.Opened {
		a := Alert{
			Timestamp: now,
			Level:     LevelAlert,
			Message:   fmt.Sprintf("Port opened: %d/%s (%s)", p.Port, p.Protocol, p.Service),
		}
		alerts = append(alerts, a)
		fmt.Fprintf(n.w, "[%s] %s %s\n", a.Timestamp.Format(time.RFC3339), a.Level, a.Message)
	}

	for _, p := range d.Closed {
		a := Alert{
			Timestamp: now,
			Level:     LevelWarn,
			Message:   fmt.Sprintf("Port closed: %d/%s (%s)", p.Port, p.Protocol, p.Service),
		}
		alerts = append(alerts, a)
		fmt.Fprintf(n.w, "[%s] %s %s\n", a.Timestamp.Format(time.RFC3339), a.Level, a.Message)
	}

	if len(alerts) == 0 {
		a := Alert{Timestamp: now, Level: LevelInfo, Message: "No port changes detected."}
		alerts = append(alerts, a)
		fmt.Fprintf(n.w, "[%s] %s %s\n", a.Timestamp.Format(time.RFC3339), a.Level, a.Message)
	}

	return alerts
}
