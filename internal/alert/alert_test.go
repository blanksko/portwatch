package alert_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makeDiff(opened, closed []scanner.Result) snapshot.DiffResult {
	return snapshot.DiffResult{Opened: opened, Closed: closed}
}

func TestNotify_OpenedPorts(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	d := makeDiff([]scanner.Result{
		{Port: 8080, Protocol: "tcp", Service: "http-alt", Open: true},
	}, nil)

	alerts := n.Notify(d)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != alert.LevelAlert {
		t.Errorf("expected level ALERT, got %s", alerts[0].Level)
	}
	if !strings.Contains(buf.String(), "8080") {
		t.Errorf("expected port 8080 in output, got: %s", buf.String())
	}
}

func TestNotify_ClosedPorts(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	d := makeDiff(nil, []scanner.Result{
		{Port: 22, Protocol: "tcp", Service: "ssh", Open: false},
	})

	alerts := n.Notify(d)
	if len(alerts) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(alerts))
	}
	if alerts[0].Level != alert.LevelWarn {
		t.Errorf("expected level WARN, got %s", alerts[0].Level)
	}
	if !strings.Contains(buf.String(), "ssh") {
		t.Errorf("expected service name in output, got: %s", buf.String())
	}
}

func TestNotify_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	n := alert.New(&buf)

	d := makeDiff(nil, nil)
	alerts := n.Notify(d)

	if len(alerts) != 1 {
		t.Fatalf("expected 1 info alert, got %d", len(alerts))
	}
	if alerts[0].Level != alert.LevelInfo {
		t.Errorf("expected level INFO, got %s", alerts[0].Level)
	}
	if !strings.Contains(buf.String(), "No port changes") {
		t.Errorf("expected no-change message, got: %s", buf.String())
	}
}

func TestNew_NilWriter(t *testing.T) {
	// Should not panic when w is nil (falls back to os.Stdout)
	n := alert.New(nil)
	if n == nil {
		t.Fatal("expected non-nil Notifier")
	}
}
