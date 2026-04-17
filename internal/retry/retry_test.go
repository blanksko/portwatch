package retry

import (
	"errors"
	"testing"
	"time"
)

var errTemp = errors.New("temporary error")

func noSleep(_ time.Duration) {}

func newFast(cfg Config) *Doer {
	d := New(cfg)
	d.sleep = noSleep
	return d
}

func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	d := newFast(Default())
	calls := 0
	err := d.Do(func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesOnFailure(t *testing.T) {
	d := newFast(Config{MaxAttempts: 3, Delay: 0})
	calls := 0
	err := d.Do(func() error {
		calls++
		if calls < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil after retries, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	d := newFast(Config{MaxAttempts: 3, Delay: 0})
	calls := 0
	err := d.Do(func() error {
		calls++
		return errTemp
	})
	if !errors.Is(err, ErrMaxAttempts) {
		t.Fatalf("expected ErrMaxAttempts, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestDo_ZeroAttempts(t *testing.T) {
	d := newFast(Config{MaxAttempts: 0, Delay: 0})
	err := d.Do(func() error { return nil })
	if !errors.Is(err, ErrMaxAttempts) {
		t.Fatalf("expected ErrMaxAttempts for zero attempts, got %v", err)
	}
}

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts 3, got %d", cfg.MaxAttempts)
	}
	if cfg.Delay != 500*time.Millisecond {
		t.Errorf("expected Delay 500ms, got %v", cfg.Delay)
	}
}
