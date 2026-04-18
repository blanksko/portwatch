package audit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/snapshot"
)

func makeDiff(opened, closed []int) snapshot.Diff {
	return snapshot.Diff{Opened: opened, Closed: closed}
}

func TestRecord_NoChanges(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	if err := l.Record("localhost", makeDiff(nil, nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var entry audit.Entry
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if entry.Message != "no changes" {
		t.Errorf("expected 'no changes', got %q", entry.Message)
	}
}

func TestRecord_OpenedPorts(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	l.Record("example.com", makeDiff([]int{80, 443}, nil))
	var entry audit.Entry
	json.Unmarshal(buf.Bytes(), &entry)
	if entry.Host != "example.com" {
		t.Errorf("expected host example.com, got %q", entry.Host)
	}
	if len(entry.Opened) != 2 {
		t.Errorf("expected 2 opened ports, got %d", len(entry.Opened))
	}
	if !strings.Contains(entry.Message, "new ports") {
		t.Errorf("unexpected message: %q", entry.Message)
	}
}

func TestRecord_ClosedPorts(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	l.Record("host", makeDiff(nil, []int{22}))
	var entry audit.Entry
	json.Unmarshal(buf.Bytes(), &entry)
	if entry.Message != "ports closed" {
		t.Errorf("unexpected message: %q", entry.Message)
	}
}

func TestRecord_BothChanges(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	l.Record("host", makeDiff([]int{8080}, []int{22}))
	var entry audit.Entry
	json.Unmarshal(buf.Bytes(), &entry)
	if entry.Message != "ports opened and closed" {
		t.Errorf("unexpected message: %q", entry.Message)
	}
}

func TestNew_NilPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil writer")
		}
	}()
	audit.New(nil)
}
