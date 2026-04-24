package trend_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/trend"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRecord_IncreasesCount(t *testing.T) {
	tr := trend.New(time.Minute)
	tr.Record("host1", 2, 0)
	tr.Record("host1", 0, 1)
	if got := tr.Count("host1"); got != 2 {
		t.Fatalf("expected 2 events, got %d", got)
	}
}

func TestCount_UnknownHostIsZero(t *testing.T) {
	tr := trend.New(time.Minute)
	if got := tr.Count("ghost"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestPrune_RemovesExpiredEvents(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	tr := trend.New(time.Minute)

	// Inject nowFunc via unexported field is not possible; use Record timing.
	// We test pruning by recording at a real past time via the public API
	// with a short window instead.
	short := trend.New(10 * time.Millisecond)
	_ = base
	short.Record("h", 1, 0)
	time.Sleep(20 * time.Millisecond)
	if got := short.Count("h"); got != 0 {
		t.Fatalf("expected 0 after window, got %d", got)
	}
}

func TestRecent_ReturnsCopy(t *testing.T) {
	tr := trend.New(time.Minute)
	tr.Record("h", 3, 1)
	events := tr.Recent("h")
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Opened != 3 || events[0].Closed != 1 {
		t.Fatalf("unexpected event values: %+v", events[0])
	}
	// Mutating the returned slice must not affect internal state.
	events[0].Opened = 99
	if tr.Recent("h")[0].Opened != 3 {
		t.Fatal("internal state was mutated through returned slice")
	}
}

func TestReset_ClearsHost(t *testing.T) {
	tr := trend.New(time.Minute)
	tr.Record("h", 1, 0)
	tr.Record("h", 2, 0)
	tr.Reset("h")
	if got := tr.Count("h"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestDifferentHosts_AreIndependent(t *testing.T) {
	tr := trend.New(time.Minute)
	tr.Record("a", 1, 0)
	tr.Record("a", 1, 0)
	tr.Record("b", 1, 0)
	if tr.Count("a") != 2 {
		t.Fatalf("host a: expected 2, got %d", tr.Count("a"))
	}
	if tr.Count("b") != 1 {
		t.Fatalf("host b: expected 1, got %d", tr.Count("b"))
	}
}
