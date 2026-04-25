package quota_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/quota"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_FirstCall(t *testing.T) {
	q := quota.New(3, time.Minute, quota.WithNow(fixedNow(time.Now())))
	if err := q.Allow("host-a"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllow_WithinLimit(t *testing.T) {
	now := time.Now()
	q := quota.New(3, time.Minute, quota.WithNow(fixedNow(now)))
	for i := 0; i < 3; i++ {
		if err := q.Allow("host-b"); err != nil {
			t.Fatalf("call %d: expected nil, got %v", i, err)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	now := time.Now()
	q := quota.New(2, time.Minute, quota.WithNow(fixedNow(now)))
	_ = q.Allow("host-c")
	_ = q.Allow("host-c")
	if err := q.Allow("host-c"); err == nil {
		t.Fatal("expected quota error, got nil")
	}
}

func TestAllow_DifferentHostsAreIndependent(t *testing.T) {
	now := time.Now()
	q := quota.New(1, time.Minute, quota.WithNow(fixedNow(now)))
	if err := q.Allow("host-x"); err != nil {
		t.Fatalf("host-x first call: %v", err)
	}
	if err := q.Allow("host-y"); err != nil {
		t.Fatalf("host-y first call: %v", err)
	}
}

func TestAllow_AfterWindowPermitted(t *testing.T) {
	base := time.Now()
	clock := base
	q := quota.New(1, time.Minute, quota.WithNow(func() time.Time { return clock }))

	_ = q.Allow("host-d")
	if err := q.Allow("host-d"); err == nil {
		t.Fatal("expected quota error before window expires")
	}

	clock = base.Add(61 * time.Second)
	if err := q.Allow("host-d"); err != nil {
		t.Fatalf("expected nil after window, got %v", err)
	}
}

func TestReset_ClearsHost(t *testing.T) {
	now := time.Now()
	q := quota.New(1, time.Minute, quota.WithNow(fixedNow(now)))
	_ = q.Allow("host-e")
	q.Reset("host-e")
	if err := q.Allow("host-e"); err != nil {
		t.Fatalf("expected nil after reset, got %v", err)
	}
}

func TestCount_ReflectsScans(t *testing.T) {
	now := time.Now()
	q := quota.New(5, time.Minute, quota.WithNow(fixedNow(now)))
	_ = q.Allow("host-f")
	_ = q.Allow("host-f")
	if got := q.Count("host-f"); got != 2 {
		t.Fatalf("expected count 2, got %d", got)
	}
}

func TestCount_UnknownHostIsZero(t *testing.T) {
	q := quota.New(5, time.Minute)
	if got := q.Count("nobody"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}
