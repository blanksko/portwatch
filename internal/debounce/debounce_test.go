package debounce

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_FirstCall(t *testing.T) {
	d := New(5 * time.Second)
	if !d.Allow("host:22") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallBlocked(t *testing.T) {
	base := time.Now()
	d := New(5 * time.Second)
	d.now = fixedNow(base)
	d.Allow("host:22")
	if d.Allow("host:22") {
		t.Fatal("expected second call within window to be blocked")
	}
}

func TestAllow_AfterWindowPermitted(t *testing.T) {
	base := time.Now()
	d := New(5 * time.Second)
	d.now = fixedNow(base)
	d.Allow("host:22")
	d.now = fixedNow(base.Add(6 * time.Second))
	if !d.Allow("host:22") {
		t.Fatal("expected call after window to be allowed")
	}
}

func TestAllow_DifferentKeysIndependent(t *testing.T) {
	d := New(5 * time.Second)
	d.Allow("host:22")
	if !d.Allow("host:80") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestReset_ClearsKey(t *testing.T) {
	base := time.Now()
	d := New(5 * time.Second)
	d.now = fixedNow(base)
	d.Allow("host:22")
	d.Reset("host:22")
	if !d.Allow("host:22") {
		t.Fatal("expected allow after reset")
	}
}

func TestPurge_RemovesExpiredKeys(t *testing.T) {
	base := time.Now()
	d := New(5 * time.Second)
	d.now = fixedNow(base)
	d.Allow("host:22")
	d.now = fixedNow(base.Add(10 * time.Second))
	d.Purge()
	if !d.Allow("host:22") {
		t.Fatal("expected key to be gone after purge")
	}
}

func TestPurge_KeepsActiveKeys(t *testing.T) {
	base := time.Now()
	d := New(5 * time.Second)
	d.now = fixedNow(base)
	d.Allow("host:22")
	d.now = fixedNow(base.Add(2 * time.Second))
	d.Purge()
	if d.Allow("host:22") {
		t.Fatal("expected active key to remain after purge")
	}
}
