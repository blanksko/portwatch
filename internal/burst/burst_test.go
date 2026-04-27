package burst_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/burst"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRecord_BelowThreshold(t *testing.T) {
	now := time.Now()
	d := burst.New(3, time.Minute, burst.WithNow(fixedNow(now)))

	for i := 0; i < 3; i++ {
		if d.Record("host-a") {
			t.Fatalf("expected no burst on event %d", i+1)
		}
	}
}

func TestRecord_AtThresholdPlusOne_ReturnsBurst(t *testing.T) {
	now := time.Now()
	d := burst.New(3, time.Minute, burst.WithNow(fixedNow(now)))

	for i := 0; i < 3; i++ {
		d.Record("host-a")
	}
	if !d.Record("host-a") {
		t.Fatal("expected burst to be detected on 4th event")
	}
}

func TestRecord_DifferentHostsAreIndependent(t *testing.T) {
	now := time.Now()
	d := burst.New(2, time.Minute, burst.WithNow(fixedNow(now)))

	d.Record("host-a")
	d.Record("host-a")
	d.Record("host-a") // bursts host-a

	if d.Record("host-b") {
		t.Fatal("host-b should not burst independently")
	}
}

func TestRecord_ExpiredEventsNotCounted(t *testing.T) {
	base := time.Now()
	clock := base
	d := burst.New(2, time.Minute, burst.WithNow(func() time.Time { return clock }))

	// record 3 events in the past (outside window)
	for i := 0; i < 3; i++ {
		d.Record("host-a")
	}
	// advance clock past the window
	clock = base.Add(2 * time.Minute)

	// fresh events should start from zero
	if d.Record("host-a") {
		t.Fatal("expired events should not contribute to burst")
	}
}

func TestCount_ReturnsActiveEventCount(t *testing.T) {
	now := time.Now()
	d := burst.New(10, time.Minute, burst.WithNow(fixedNow(now)))

	d.Record("host-a")
	d.Record("host-a")

	if got := d.Count("host-a"); got != 2 {
		t.Fatalf("expected count 2, got %d", got)
	}
}

func TestCount_UnknownHostIsZero(t *testing.T) {
	d := burst.New(3, time.Minute)
	if got := d.Count("ghost"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestReset_ClearsHost(t *testing.T) {
	now := time.Now()
	d := burst.New(2, time.Minute, burst.WithNow(fixedNow(now)))

	d.Record("host-a")
	d.Record("host-a")
	d.Record("host-a")
	d.Reset("host-a")

	if got := d.Count("host-a"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}
