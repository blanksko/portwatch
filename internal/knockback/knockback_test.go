package knockback

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func newGate(now func() time.Time) *Gate {
	g := New()
	g.now = now
	return g
}

func TestAllow_NewHostPermitted(t *testing.T) {
	g := New()
	if !g.Allow("host1") {
		t.Fatal("expected new host to be allowed")
	}
}

func TestAllow_BelowThresholdPermitted(t *testing.T) {
	g := New()
	g.RecordFailure("h")
	g.RecordFailure("h")
	if !g.Allow("h") {
		t.Fatal("expected host below threshold to be allowed")
	}
}

func TestAllow_AtThresholdBlocked(t *testing.T) {
	now := time.Now()
	g := newGate(fixedNow(now))
	for i := 0; i < DefaultThreshold; i++ {
		g.RecordFailure("h")
	}
	if g.Allow("h") {
		t.Fatal("expected host at threshold to be blocked")
	}
}

func TestAllow_AfterBackoffPermitted(t *testing.T) {
	now := time.Now()
	g := newGate(fixedNow(now))
	for i := 0; i < DefaultThreshold; i++ {
		g.RecordFailure("h")
	}
	// Advance clock past the base back-off window.
	g.now = fixedNow(now.Add(DefaultBase + time.Second))
	if !g.Allow("h") {
		t.Fatal("expected host to be allowed after back-off window")
	}
}

func TestRecordSuccess_ResetsState(t *testing.T) {
	now := time.Now()
	g := newGate(fixedNow(now))
	for i := 0; i < DefaultThreshold; i++ {
		g.RecordFailure("h")
	}
	g.RecordSuccess("h")
	if !g.Allow("h") {
		t.Fatal("expected host to be allowed after success reset")
	}
	if g.Failures("h") != 0 {
		t.Fatalf("expected 0 failures after reset, got %d", g.Failures("h"))
	}
}

func TestBackoff_Doubles(t *testing.T) {
	now := time.Now()
	g := newGate(fixedNow(now))
	g.base = 10 * time.Second
	g.max = 5 * time.Minute
	for i := 0; i < DefaultThreshold; i++ {
		g.RecordFailure("h")
	}
	// Advance past first window, trigger second block.
	g.now = fixedNow(now.Add(g.base + time.Second))
	g.RecordFailure("h")
	e := g.state["h"]
	if e.backoff != 20*time.Second {
		t.Fatalf("expected doubled backoff 20s, got %v", e.backoff)
	}
}

func TestBackoff_CapsAtMax(t *testing.T) {
	now := time.Now()
	g := newGate(fixedNow(now))
	g.base = 1 * time.Second
	g.max = 4 * time.Second
	for i := 0; i < 10; i++ {
		g.RecordFailure("h")
		g.now = fixedNow(g.now().Add(g.max + time.Second))
	}
	e := g.state["h"]
	if e.backoff > g.max {
		t.Fatalf("backoff %v exceeds max %v", e.backoff, g.max)
	}
}

func TestFailures_UnknownHostIsZero(t *testing.T) {
	g := New()
	if n := g.Failures("unknown"); n != 0 {
		t.Fatalf("expected 0, got %d", n)
	}
}
