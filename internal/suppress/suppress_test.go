package suppress

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_FirstCall(t *testing.T) {
	s := New(time.Minute)
	if !s.Allow("localhost") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallBlocked(t *testing.T) {
	now := time.Now()
	s := New(time.Minute)
	s.now = fixedNow(now)
	s.Allow("localhost")
	if s.Allow("localhost") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_AfterCooldown(t *testing.T) {
	now := time.Now()
	s := New(time.Minute)
	s.now = fixedNow(now)
	s.Allow("localhost")
	s.now = fixedNow(now.Add(2 * time.Minute))
	if !s.Allow("localhost") {
		t.Fatal("expected call after cooldown to be allowed")
	}
}

func TestAllow_DifferentHostsIndependent(t *testing.T) {
	s := New(time.Minute)
	s.Allow("host-a")
	if !s.Allow("host-b") {
		t.Fatal("expected different host to be allowed")
	}
}

func TestReset_ClearsHost(t *testing.T) {
	now := time.Now()
	s := New(time.Minute)
	s.now = fixedNow(now)
	s.Allow("localhost")
	s.Reset("localhost")
	if !s.Allow("localhost") {
		t.Fatal("expected allow after reset")
	}
}

func TestResetAll_ClearsAll(t *testing.T) {
	now := time.Now()
	s := New(time.Minute)
	s.now = fixedNow(now)
	s.Allow("host-a")
	s.Allow("host-b")
	s.ResetAll()
	if !s.Allow("host-a") || !s.Allow("host-b") {
		t.Fatal("expected all hosts allowed after ResetAll")
	}
}
