package backoff

import (
	"testing"
	"time"
)

func TestNext_FirstCallReturnsBase(t *testing.T) {
	b := New(1*time.Second, 60*time.Second, 2.0, 0)
	d := b.Next("host1")
	// with zero jitter the first call should equal base exactly
	if d != 1*time.Second {
		t.Fatalf("expected 1s, got %v", d)
	}
}

func TestNext_SecondCallDoubles(t *testing.T) {
	b := New(1*time.Second, 60*time.Second, 2.0, 0)
	_ = b.Next("host1")
	d := b.Next("host1")
	if d != 2*time.Second {
		t.Fatalf("expected 2s, got %v", d)
	}
}

func TestNext_CapsAtMax(t *testing.T) {
	b := New(1*time.Second, 4*time.Second, 2.0, 0)
	for i := 0; i < 10; i++ {
		_ = b.Next("host1")
	}
	d := b.Next("host1")
	if d > 4*time.Second {
		t.Fatalf("expected <= 4s, got %v", d)
	}
}

func TestReset_ClearsCount(t *testing.T) {
	b := New(1*time.Second, 60*time.Second, 2.0, 0)
	_ = b.Next("host1")
	_ = b.Next("host1")
	b.Reset("host1")
	if b.Failures("host1") != 0 {
		t.Fatal("expected failures to be 0 after reset")
	}
	d := b.Next("host1")
	if d != 1*time.Second {
		t.Fatalf("expected 1s after reset, got %v", d)
	}
}

func TestFailures_TracksCount(t *testing.T) {
	b := Default()
	if b.Failures("h") != 0 {
		t.Fatal("expected 0 initial failures")
	}
	_ = b.Next("h")
	_ = b.Next("h")
	if b.Failures("h") != 2 {
		t.Fatalf("expected 2 failures, got %d", b.Failures("h"))
	}
}

func TestNext_DifferentHostsAreIndependent(t *testing.T) {
	b := New(1*time.Second, 60*time.Second, 2.0, 0)
	_ = b.Next("a")
	_ = b.Next("a")
	_ = b.Next("a")
	d := b.Next("b")
	if d != 1*time.Second {
		t.Fatalf("host b should start fresh, got %v", d)
	}
}

func TestDefault_ReturnsNonNil(t *testing.T) {
	b := Default()
	if b == nil {
		t.Fatal("expected non-nil backoff")
	}
	d := b.Next("any")
	if d <= 0 {
		t.Fatalf("expected positive duration, got %v", d)
	}
}
