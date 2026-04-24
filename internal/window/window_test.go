package window

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestCount_EmptyKey(t *testing.T) {
	c := New(time.Minute)
	if got := c.Count("host-a"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_IncrementsCount(t *testing.T) {
	c := New(time.Minute)
	c.Add("host-a")
	c.Add("host-a")
	if got := c.Count("host-a"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestCount_PrunesExpiredEvents(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	c := New(time.Minute)

	c.now = fixedNow(base)
	c.Add("host-a")
	c.Add("host-a")

	// Advance time beyond the window.
	c.now = fixedNow(base.Add(90 * time.Second))
	c.Add("host-a") // one fresh event

	if got := c.Count("host-a"); got != 1 {
		t.Fatalf("expected 1 after pruning, got %d", got)
	}
}

func TestReset_ClearsKey(t *testing.T) {
	c := New(time.Minute)
	c.Add("host-a")
	c.Add("host-a")
	c.Reset("host-a")
	if got := c.Count("host-a"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestCount_DifferentKeysIndependent(t *testing.T) {
	c := New(time.Minute)
	c.Add("host-a")
	c.Add("host-a")
	c.Add("host-b")

	if got := c.Count("host-a"); got != 2 {
		t.Fatalf("host-a: expected 2, got %d", got)
	}
	if got := c.Count("host-b"); got != 1 {
		t.Fatalf("host-b: expected 1, got %d", got)
	}
}

func TestNew_ZeroDurationPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero duration")
		}
	}()
	New(0)
}

func TestCount_EventsAtWindowBoundaryExcluded(t *testing.T) {
	base := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	c := New(time.Minute)

	c.now = fixedNow(base)
	c.Add("host-a")

	// Advance exactly to the boundary — event is now expired.
	c.now = fixedNow(base.Add(time.Minute))

	if got := c.Count("host-a"); got != 0 {
		t.Fatalf("expected 0 at exact boundary, got %d", got)
	}
}
