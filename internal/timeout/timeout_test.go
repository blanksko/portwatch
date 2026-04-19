package timeout_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/timeout"
)

func TestNew_DefaultTimeout(t *testing.T) {
	m := timeout.New(3 * time.Second)
	if m.Default() != 3*time.Second {
		t.Fatalf("expected 3s, got %v", m.Default())
	}
}

func TestNew_ZeroUsesDefault(t *testing.T) {
	m := timeout.New(0)
	if m.Default() != 5*time.Second {
		t.Fatalf("expected fallback 5s, got %v", m.Default())
	}
}

func TestGet_FallsBackToDefault(t *testing.T) {
	m := timeout.New(2 * time.Second)
	if got := m.Get("192.168.1.1"); got != 2*time.Second {
		t.Fatalf("expected 2s, got %v", got)
	}
}

func TestSet_OverridesDefault(t *testing.T) {
	m := timeout.New(2 * time.Second)
	m.Set("slow-host", 10*time.Second)
	if got := m.Get("slow-host"); got != 10*time.Second {
		t.Fatalf("expected 10s, got %v", got)
	}
}

func TestDelete_RevertsToDefault(t *testing.T) {
	m := timeout.New(2 * time.Second)
	m.Set("slow-host", 10*time.Second)
	m.Delete("slow-host")
	if got := m.Get("slow-host"); got != 2*time.Second {
		t.Fatalf("expected default 2s after delete, got %v", got)
	}
}

func TestSet_MultipleHosts(t *testing.T) {
	m := timeout.New(1 * time.Second)
	m.Set("host-a", 4*time.Second)
	m.Set("host-b", 8*time.Second)

	if got := m.Get("host-a"); got != 4*time.Second {
		t.Errorf("host-a: expected 4s, got %v", got)
	}
	if got := m.Get("host-b"); got != 8*time.Second {
		t.Errorf("host-b: expected 8s, got %v", got)
	}
	if got := m.Get("host-c"); got != 1*time.Second {
		t.Errorf("host-c: expected default 1s, got %v", got)
	}
}
