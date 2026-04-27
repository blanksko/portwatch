package probe_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/probe"
)

func TestNew_DefaultInterval(t *testing.T) {
	m := probe.New(0)
	if m.Default() != 60*time.Second {
		t.Fatalf("expected 60s default, got %s", m.Default())
	}
}

func TestNew_CustomDefault(t *testing.T) {
	m := probe.New(30 * time.Second)
	if m.Default() != 30*time.Second {
		t.Fatalf("expected 30s, got %s", m.Default())
	}
}

func TestGet_FallsBackToDefault(t *testing.T) {
	m := probe.New(45 * time.Second)
	if got := m.Get("192.168.1.1"); got != 45*time.Second {
		t.Fatalf("expected default 45s, got %s", got)
	}
}

func TestSet_OverridesDefault(t *testing.T) {
	m := probe.New(60 * time.Second)
	if err := m.Set("10.0.0.1", 10*time.Second); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := m.Get("10.0.0.1"); got != 10*time.Second {
		t.Fatalf("expected 10s, got %s", got)
	}
}

func TestSet_InvalidInterval(t *testing.T) {
	m := probe.New(60 * time.Second)
	if err := m.Set("10.0.0.1", 0); err == nil {
		t.Fatal("expected error for zero interval")
	}
	if err := m.Set("10.0.0.1", -5*time.Second); err == nil {
		t.Fatal("expected error for negative interval")
	}
}

func TestDelete_RevertsToDefault(t *testing.T) {
	m := probe.New(60 * time.Second)
	_ = m.Set("host-a", 5*time.Second)
	m.Delete("host-a")
	if got := m.Get("host-a"); got != 60*time.Second {
		t.Fatalf("expected default after delete, got %s", got)
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	m := probe.New(60 * time.Second)
	_ = m.Set("a", 10*time.Second)
	_ = m.Set("b", 20*time.Second)
	all := m.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if all["a"] != 10*time.Second {
		t.Errorf("wrong interval for a: %s", all["a"])
	}
	if all["b"] != 20*time.Second {
		t.Errorf("wrong interval for b: %s", all["b"])
	}
}

func TestAll_MutationDoesNotAffectInternal(t *testing.T) {
	m := probe.New(60 * time.Second)
	_ = m.Set("x", 15*time.Second)
	all := m.All()
	all["x"] = 999 * time.Second
	if got := m.Get("x"); got != 15*time.Second {
		t.Fatalf("internal state was mutated: %s", got)
	}
}
