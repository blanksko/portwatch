package throttle

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_FirstCall(t *testing.T) {
	th := New(time.Minute)
	if !th.Allow("localhost") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallBlocked(t *testing.T) {
	base := time.Now()
	th := New(time.Minute)
	th.now = fixedNow(base)
	th.Allow("localhost")
	if th.Allow("localhost") {
		t.Fatal("expected second immediate call to be blocked")
	}
}

func TestAllow_AfterIntervalPermitted(t *testing.T) {
	base := time.Now()
	th := New(time.Minute)
	th.now = fixedNow(base)
	th.Allow("localhost")
	th.now = fixedNow(base.Add(2 * time.Minute))
	if !th.Allow("localhost") {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestAllow_DifferentHostsIndependent(t *testing.T) {
	base := time.Now()
	th := New(time.Minute)
	th.now = fixedNow(base)
	th.Allow("host-a")
	if !th.Allow("host-b") {
		t.Fatal("expected different host to be allowed")
	}
}

func TestReset_AllowsImmediateScan(t *testing.T) {
	base := time.Now()
	th := New(time.Minute)
	th.now = fixedNow(base)
	th.Allow("localhost")
	th.Reset("localhost")
	if !th.Allow("localhost") {
		t.Fatal("expected allow after reset")
	}
}

func TestNextAllowed_NeverScanned(t *testing.T) {
	th := New(time.Minute)
	if !th.NextAllowed("localhost").IsZero() {
		t.Fatal("expected zero time for unscanned host")
	}
}

func TestNextAllowed_AfterScan(t *testing.T) {
	base := time.Now()
	th := New(time.Minute)
	th.now = fixedNow(base)
	th.Allow("localhost")
	want := base.Add(time.Minute)
	if got := th.NextAllowed("localhost"); !got.Equal(want) {
		t.Fatalf("NextAllowed = %v, want %v", got, want)
	}
}
