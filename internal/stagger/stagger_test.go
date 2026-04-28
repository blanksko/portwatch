package stagger_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/stagger"
)

func TestDelay_AssignsStableOffset(t *testing.T) {
	s := stagger.New(10*time.Second, 10)

	first := s.Delay("host-a")
	second := s.Delay("host-a")

	if first != second {
		t.Fatalf("expected stable delay for same host, got %v then %v", first, second)
	}
}

func TestDelay_DifferentHostsGetDifferentSlots(t *testing.T) {
	s := stagger.New(10*time.Second, 10)

	a := s.Delay("host-a")
	b := s.Delay("host-b")

	if a == b {
		t.Fatalf("expected different delays, both got %v", a)
	}
}

func TestDelay_WithinWindow(t *testing.T) {
	window := 5 * time.Second
	s := stagger.New(window, 5)

	hosts := []string{"h1", "h2", "h3", "h4", "h5"}
	for _, h := range hosts {
		d := s.Delay(h)
		if d < 0 || d >= window {
			t.Errorf("delay %v for %s is outside window [0, %v)", d, h, window)
		}
	}
}

func TestReset_ClearsSlots(t *testing.T) {
	s := stagger.New(10*time.Second, 10)

	before := s.Delay("host-a")
	s.Reset()
	after := s.Delay("host-a")

	// After reset the host gets slot 0 again.
	if after != 0 {
		t.Fatalf("expected slot 0 after reset, got %v", after)
	}
	_ = before
}

func TestLen_TracksHosts(t *testing.T) {
	s := stagger.New(10*time.Second, 10)

	if s.Len() != 0 {
		t.Fatal("expected Len 0 initially")
	}

	s.Delay("h1")
	s.Delay("h2")
	s.Delay("h1") // duplicate – should not increase count

	if s.Len() != 2 {
		t.Fatalf("expected Len 2, got %d", s.Len())
	}
}

func TestNew_ZeroNDefaultsToOne(t *testing.T) {
	s := stagger.New(10*time.Second, 0)
	d := s.Delay("host")
	if d != 0 {
		t.Fatalf("expected delay 0 for single slot, got %v", d)
	}
}
