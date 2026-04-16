package notifier_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/notifier"
	"github.com/user/portwatch/internal/snapshot"
)

func makeDiff(opened, closed []int) snapshot.DiffResult {
	return snapshot.DiffResult{Opened: opened, Closed: closed}
}

func TestNotify_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	n, _ := notifier.New(&buf)
	if err := n.Notify("localhost", makeDiff(nil, nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected no output, got: %q", buf.String())
	}
}

func TestNotify_OpenedPorts(t *testing.T) {
	var buf bytes.Buffer
	n, _ := notifier.New(&buf)
	_ = n.Notify("myhost", makeDiff([]int{80, 443}, nil))
	out := buf.String()
	if !strings.Contains(out, "myhost") {
		t.Errorf("expected host in output")
	}
	if !strings.Contains(out, "+ 80") || !strings.Contains(out, "+ 443") {
		t.Errorf("expected opened ports in output, got: %q", out)
	}
}

func TestNotify_ClosedPorts(t *testing.T) {
	var buf bytes.Buffer
	n, _ := notifier.New(&buf)
	_ = n.Notify("myhost", makeDiff(nil, []int{22}))
	out := buf.String()
	if !strings.Contains(out, "- 22") {
		t.Errorf("expected closed port in output, got: %q", out)
	}
}

func TestNotify_BothChanges(t *testing.T) {
	var buf bytes.Buffer
	n, _ := notifier.New(&buf)
	_ = n.Notify("host", makeDiff([]int{8080}, []int{3306}))
	out := buf.String()
	if !strings.Contains(out, "+ 8080") || !strings.Contains(out, "- 3306") {
		t.Errorf("expected both changes in output, got: %q", out)
	}
}

func TestNew_NilWriter(t *testing.T) {
	_, err := notifier.New(nil)
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}
