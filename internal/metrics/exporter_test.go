package metrics

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func TestExport_ContainsMetrics(t *testing.T) {
	m := New()
	results := []scanner.Result{
		{Host: "localhost", Port: 80, Open: true},
		{Host: "localhost", Port: 81, Open: false},
	}
	m.Record(results)

	e := NewExporter(m)
	var buf bytes.Buffer
	if err := e.Export(&buf); err != nil {
		t.Fatalf("Export returned error: %v", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(buf.Bytes(), &snap); err != nil {
		t.Fatalf("failed to unmarshal snapshot: %v", err)
	}

	if snap.TotalScans != 1 {
		t.Errorf("expected TotalScans=1, got %d", snap.TotalScans)
	}
	if snap.OpenPorts != 1 {
		t.Errorf("expected OpenPorts=1, got %d", snap.OpenPorts)
	}
	if snap.ClosedPorts != 1 {
		t.Errorf("expected ClosedPorts=1, got %d", snap.ClosedPorts)
	}
	if snap.ExportedAt.IsZero() {
		t.Error("expected ExportedAt to be set")
	}
}

func TestExport_EmptyMetrics(t *testing.T) {
	m := New()
	e := NewExporter(m)
	var buf bytes.Buffer
	if err := e.Export(&buf); err != nil {
		t.Fatalf("Export returned error: %v", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(buf.Bytes(), &snap); err != nil {
		t.Fatalf("failed to unmarshal snapshot: %v", err)
	}
	if snap.TotalScans != 0 {
		t.Errorf("expected TotalScans=0, got %d", snap.TotalScans)
	}
}

func TestExport_LastScanTime(t *testing.T) {
	m := New()
	before := time.Now()
	m.Record([]scanner.Result{{Host: "h", Port: 22, Open: true}})
	after := time.Now()

	e := NewExporter(m)
	var buf bytes.Buffer
	_ = e.Export(&buf)

	var snap Snapshot
	_ = json.Unmarshal(buf.Bytes(), &snap)

	if snap.LastScanTime.Before(before) || snap.LastScanTime.After(after) {
		t.Errorf("LastScanTime %v not in expected range", snap.LastScanTime)
	}
}

func TestNewExporter_NilPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil Metrics")
		}
	}()
	NewExporter(nil)
}
