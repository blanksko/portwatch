package labels_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/labels"
)

func TestGet_KnownHost(t *testing.T) {
	l := labels.New(map[string]string{"192.168.1.1": "gateway"})
	v, ok := l.Get("192.168.1.1")
	if !ok || v != "gateway" {
		t.Fatalf("expected 'gateway', got %q ok=%v", v, ok)
	}
}

func TestGet_UnknownHost(t *testing.T) {
	l := labels.New(nil)
	_, ok := l.Get("10.0.0.1")
	if ok {
		t.Fatal("expected no label for unknown host")
	}
}

func TestSet_AddsLabel(t *testing.T) {
	l := labels.New(nil)
	l.Set("10.0.0.2", "db-server")
	v, ok := l.Get("10.0.0.2")
	if !ok || v != "db-server" {
		t.Fatalf("expected 'db-server', got %q", v)
	}
}

func TestSet_OverwritesLabel(t *testing.T) {
	l := labels.New(map[string]string{"host": "old"})
	l.Set("host", "new")
	v, _ := l.Get("host")
	if v != "new" {
		t.Fatalf("expected 'new', got %q", v)
	}
}

func TestDelete_RemovesLabel(t *testing.T) {
	l := labels.New(map[string]string{"host": "label"})
	l.Delete("host")
	_, ok := l.Get("host")
	if ok {
		t.Fatal("expected label to be deleted")
	}
}

// TestDelete_NonexistentKey ensures Delete is a no-op for keys not in the store.
func TestDelete_NonexistentKey(t *testing.T) {
	l := labels.New(nil)
	// Should not panic or error.
	l.Delete("ghost")
	if _, ok := l.Get("ghost"); ok {
		t.Fatal("expected no entry for never-set key")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	initial := map[string]string{"a": "alpha", "b": "beta"}
	l := labels.New(initial)
	all := l.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// mutating the returned map must not affect the store
	all["c"] = "gamma"
	if _, ok := l.Get("c"); ok {
		t.Fatal("mutation of All() result leaked into store")
	}
}

func TestNew_NilMap(t *testing.T) {
	l := labels.New(nil)
	if got := l.All(); len(got) != 0 {
		t.Fatalf("expected empty store, got %v", got)
	}
}
