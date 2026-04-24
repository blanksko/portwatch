package jitter_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/jitter"
)

func TestApply_ZeroBase(t *testing.T) {
	j := jitter.New(0.25)
	if got := j.Apply(0); got != 0 {
		t.Fatalf("expected 0, got %v", got)
	}
}

func TestApply_ZeroFraction(t *testing.T) {
	j := jitter.New(0)
	base := 10 * time.Second
	if got := j.Apply(base); got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}

func TestApply_WithinBounds(t *testing.T) {
	const frac = 0.25
	j := jitter.New(frac)
	base := 100 * time.Millisecond
	lo := time.Duration(float64(base) * (1 - frac))
	hi := time.Duration(float64(base) * (1 + frac))

	for i := 0; i < 500; i++ {
		got := j.Apply(base)
		if got < lo || got > hi {
			t.Fatalf("iteration %d: %v not in [%v, %v]", i, got, lo, hi)
		}
	}
}

func TestApply_NegativeBaseReturnsZero(t *testing.T) {
	j := jitter.New(0.5)
	// Negative base should not produce a negative result.
	if got := j.Apply(-1 * time.Second); got < 0 {
		t.Fatalf("expected >= 0, got %v", got)
	}
}

func TestNew_ClampsFractionAboveOne(t *testing.T) {
	j := jitter.New(5.0)
	base := 100 * time.Millisecond
	for i := 0; i < 200; i++ {
		if got := j.Apply(base); got < 0 {
			t.Fatalf("negative duration %v", got)
		}
	}
}

func TestNew_ClampsFractionBelowZero(t *testing.T) {
	j := jitter.New(-1.0)
	base := 50 * time.Millisecond
	if got := j.Apply(base); got != base {
		t.Fatalf("expected %v unchanged, got %v", base, got)
	}
}

func TestApply_ProducesVariance(t *testing.T) {
	j := jitter.New(0.3)
	base := 200 * time.Millisecond
	seen := make(map[time.Duration]struct{})
	for i := 0; i < 50; i++ {
		seen[j.Apply(base)] = struct{}{}
	}
	if len(seen) < 5 {
		t.Fatalf("expected variance across 50 calls, got only %d distinct values", len(seen))
	}
}
