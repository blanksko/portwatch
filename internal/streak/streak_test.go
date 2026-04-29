package streak

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestCount_UnknownHostIsZero(t *testing.T) {
	s := New()
	if got := s.Count("host-a"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestRecord_NoChange_Increments(t *testing.T) {
	s := New()
	s.Record("host-a", false)
	s.Record("host-a", false)
	s.Record("host-a", false)
	if got := s.Count("host-a"); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestRecord_Changed_ResetsToZero(t *testing.T) {
	s := New()
	s.Record("host-a", false)
	s.Record("host-a", false)
	s.Record("host-a", true)
	if got := s.Count("host-a"); got != 0 {
		t.Fatalf("expected 0 after change, got %d", got)
	}
}

func TestRecord_AfterReset_Resumes(t *testing.T) {
	s := New()
	s.Record("host-a", false)
	s.Record("host-a", false)
	s.Record("host-a", true) // reset
	s.Record("host-a", false)
	s.Record("host-a", false)
	if got := s.Count("host-a"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestRecord_DifferentHostsAreIndependent(t *testing.T) {
	s := New()
	s.Record("host-a", false)
	s.Record("host-a", false)
	s.Record("host-b", false)
	if s.Count("host-a") != 2 {
		t.Fatalf("host-a: expected 2")
	}
	if s.Count("host-b") != 1 {
		t.Fatalf("host-b: expected 1")
	}
}

func TestLastSeen_UnknownHost(t *testing.T) {
	s := New()
	_, ok := s.LastSeen("host-z")
	if ok {
		t.Fatal("expected not-ok for unknown host")
	}
}

func TestLastSeen_KnownHost(t *testing.T) {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	s := New()
	s.now = fixedNow(now)
	s.Record("host-a", false)
	got, ok := s.LastSeen("host-a")
	if !ok {
		t.Fatal("expected ok")
	}
	if !got.Equal(now) {
		t.Fatalf("expected %v, got %v", now, got)
	}
}

func TestReset_ClearsHost(t *testing.T) {
	s := New()
	s.Record("host-a", false)
	s.Record("host-a", false)
	s.Reset("host-a")
	if got := s.Count("host-a"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
	_, ok := s.LastSeen("host-a")
	if ok {
		t.Fatal("expected last-seen to be cleared")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := New()
	s.Record("host-a", false)
	s.Record("host-a", false)
	s.Record("host-b", false)
	all := s.All()
	if all["host-a"] != 2 || all["host-b"] != 1 {
		t.Fatalf("unexpected snapshot: %v", all)
	}
	// mutating the snapshot must not affect internal state
	all["host-a"] = 99
	if s.Count("host-a") != 2 {
		t.Fatal("internal state was mutated through snapshot")
	}
}
