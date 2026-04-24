package cooldown_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/cooldown"
)

func fixedNow(t time.Time) cooldown.Clock {
	return func() time.Time { return t }
}

func TestAllow_FirstCall(t *testing.T) {
	cd := cooldown.New(time.Minute, fixedNow(time.Now()))
	if !cd.Allow("host1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallBlocked(t *testing.T) {
	now := time.Now()
	cd := cooldown.New(time.Minute, fixedNow(now))
	cd.Allow("host1")
	if cd.Allow("host1") {
		t.Fatal("expected second call within window to be blocked")
	}
}

func TestAllow_AfterWindowPermitted(t *testing.T) {
	base := time.Now()
	calls := 0
	clock := func() time.Time {
		calls++
		if calls <= 2 {
			return base
		}
		return base.Add(2 * time.Minute)
	}
	cd := cooldown.New(time.Minute, clock)
	cd.Allow("host1")
	cd.Allow("host1") // blocked
	if !cd.Allow("host1") {
		t.Fatal("expected call after window to be allowed")
	}
}

func TestAllow_DifferentHostsIndependent(t *testing.T) {
	now := time.Now()
	cd := cooldown.New(time.Minute, fixedNow(now))
	cd.Allow("host1")
	if !cd.Allow("host2") {
		t.Fatal("expected different host to be allowed independently")
	}
}

func TestReset_ClearsHost(t *testing.T) {
	now := time.Now()
	cd := cooldown.New(time.Minute, fixedNow(now))
	cd.Allow("host1")
	cd.Reset("host1")
	if !cd.Allow("host1") {
		t.Fatal("expected host to be allowed after reset")
	}
}

func TestRemaining_ZeroWhenNotSeen(t *testing.T) {
	cd := cooldown.New(time.Minute, time.Now)
	if r := cd.Remaining("unknown"); r != 0 {
		t.Fatalf("expected 0 remaining for unknown host, got %v", r)
	}
}

func TestRemaining_PositiveWithinWindow(t *testing.T) {
	now := time.Now()
	cd := cooldown.New(time.Minute, fixedNow(now))
	cd.Allow("host1")
	if r := cd.Remaining("host1"); r != time.Minute {
		t.Fatalf("expected 1m remaining, got %v", r)
	}
}

func TestNew_DefaultWindow(t *testing.T) {
	cd := cooldown.New(0, time.Now)
	if cd == nil {
		t.Fatal("expected non-nil cooldown with zero window")
	}
}
