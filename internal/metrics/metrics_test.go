package metrics_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestRecord_UpdatesSnapshot(t *testing.T) {
	c := metrics.New()
	c.Record(5, 20*time.Millisecond)

	s := c.Snapshot()
	if s.TotalScans != 1 {
		t.Errorf("expected TotalScans=1, got %d", s.TotalScans)
	}
	if s.OpenPorts != 5 {
		t.Errorf("expected OpenPorts=5, got %d", s.OpenPorts)
	}
	if s.LastDuration != 20*time.Millisecond {
		t.Errorf("unexpected duration %v", s.LastDuration)
	}
	if s.LastScanAt.IsZero() {
		t.Error("expected LastScanAt to be set")
	}
}

func TestRecord_AccumulatesScans(t *testing.T) {
	c := metrics.New()
	c.Record(3, 10*time.Millisecond)
	c.Record(7, 15*time.Millisecond)

	s := c.Snapshot()
	if s.TotalScans != 2 {
		t.Errorf("expected TotalScans=2, got %d", s.TotalScans)
	}
	// OpenPorts reflects last scan only
	if s.OpenPorts != 7 {
		t.Errorf("expected OpenPorts=7, got %d", s.OpenPorts)
	}
}

func TestReset_ClearsMetrics(t *testing.T) {
	c := metrics.New()
	c.Record(4, 5*time.Millisecond)
	c.Reset()

	s := c.Snapshot()
	if s.TotalScans != 0 {
		t.Errorf("expected TotalScans=0 after reset, got %d", s.TotalScans)
	}
	if s.OpenPorts != 0 {
		t.Errorf("expected OpenPorts=0 after reset, got %d", s.OpenPorts)
	}
}

func TestNew_InitialState(t *testing.T) {
	c := metrics.New()
	s := c.Snapshot()
	if s.TotalScans != 0 || s.OpenPorts != 0 {
		t.Error("expected zero-value initial snapshot")
	}
	if !s.LastScanAt.IsZero() {
		t.Error("expected LastScanAt to be zero initially")
	}
}
