package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestAllow_FirstCallPermitted(t *testing.T) {
	l := ratelimit.New(time.Minute)
	if !l.Allow("host:80") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallBlocked(t *testing.T) {
	l := ratelimit.New(time.Minute)
	l.Allow("host:80")
	if l.Allow("host:80") {
		t.Fatal("expected second call within interval to be blocked")
	}
}

func TestAllow_DifferentKeysIndependent(t *testing.T) {
	l := ratelimit.New(time.Minute)
	l.Allow("host:80")
	if !l.Allow("host:443") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestAllow_AfterIntervalPermitted(t *testing.T) {
	l := ratelimit.New(10 * time.Millisecond)
	l.Allow("host:80")
	time.Sleep(20 * time.Millisecond)
	if !l.Allow("host:80") {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestReset_ClearsKey(t *testing.T) {
	l := ratelimit.New(time.Minute)
	l.Allow("host:80")
	l.Reset("host:80")
	if !l.Allow("host:80") {
		t.Fatal("expected allow after reset")
	}
}

func TestResetAll_ClearsAllKeys(t *testing.T) {
	l := ratelimit.New(time.Minute)
	l.Allow("host:80")
	l.Allow("host:443")
	l.ResetAll()
	if !l.Allow("host:80") || !l.Allow("host:443") {
		t.Fatal("expected all keys to be cleared after ResetAll")
	}
}

func TestNew_ZeroInterval_UsesDefault(t *testing.T) {
	l := ratelimit.New(0)
	if l == nil {
		t.Fatal("expected non-nil limiter")
	}
	if !l.Allow("key") {
		t.Fatal("expected first call to be allowed")
	}
}
