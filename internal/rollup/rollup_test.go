package rollup

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/snapshot"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAdd_WithinWindow_NoEmit(t *testing.T) {
	r := New(5 * time.Second)
	r.now = fixedNow(epoch)

	d := snapshot.Diff{Opened: []int{80}}
	_, ok := r.Add(d)
	if ok {
		t.Fatal("expected no emit within window")
	}
}

func TestAdd_AfterWindow_Emits(t *testing.T) {
	r := New(5 * time.Second)
	r.now = fixedNow(epoch)
	r.Add(snapshot.Diff{Opened: []int{80}})

	r.now = fixedNow(epoch.Add(6 * time.Second))
	merged, ok := r.Add(snapshot.Diff{Opened: []int{443}})
	if !ok {
		t.Fatal("expected emit after window")
	}
	if len(merged.Opened) != 2 {
		t.Fatalf("expected 2 opened ports, got %d", len(merged.Opened))
	}
}

func TestAdd_OpenThenClose_CancelsOut(t *testing.T) {
	r := New(5 * time.Second)
	r.now = fixedNow(epoch)
	r.Add(snapshot.Diff{Opened: []int{8080}})

	r.now = fixedNow(epoch.Add(1 * time.Second))
	r.Add(snapshot.Diff{Closed: []int{8080}})

	r.now = fixedNow(epoch.Add(6 * time.Second))
	merged, ok := r.Add(snapshot.Diff{})
	if !ok {
		t.Fatal("expected emit")
	}
	if len(merged.Opened) != 0 || len(merged.Closed) != 0 {
		t.Fatalf("expected empty diff after cancel-out, got opened=%v closed=%v",
			merged.Opened, merged.Closed)
	}
}

func TestFlush_ReturnsAccumulated(t *testing.T) {
	r := New(1 * time.Minute)
	r.now = fixedNow(epoch)
	r.Add(snapshot.Diff{Opened: []int{22, 80}})

	merged := r.Flush()
	if len(merged.Opened) != 2 {
		t.Fatalf("expected 2 opened ports after flush, got %d", len(merged.Opened))
	}
}

func TestFlush_ResetsState(t *testing.T) {
	r := New(1 * time.Minute)
	r.now = fixedNow(epoch)
	r.Add(snapshot.Diff{Opened: []int{22}})
	r.Flush()

	second := r.Flush()
	if len(second.Opened) != 0 {
		t.Fatalf("expected empty diff after second flush, got %v", second.Opened)
	}
}

func TestAdd_ResetsWindowAfterEmit(t *testing.T) {
	r := New(5 * time.Second)
	r.now = fixedNow(epoch)
	r.Add(snapshot.Diff{Opened: []int{80}})

	r.now = fixedNow(epoch.Add(6 * time.Second))
	r.Add(snapshot.Diff{})

	// New window starts; should not emit immediately.
	r.now = fixedNow(epoch.Add(7 * time.Second))
	_, ok := r.Add(snapshot.Diff{Opened: []int{443}})
	if ok {
		t.Fatal("expected no emit within new window")
	}
}
