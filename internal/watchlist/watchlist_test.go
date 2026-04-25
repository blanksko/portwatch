package watchlist_test

import (
	"testing"

	"github.com/user/portwatch/internal/watchlist"
)

func TestAdd_NewHost(t *testing.T) {
	w := watchlist.New()
	if err := w.Add("192.168.1.1", "router"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !w.Has("192.168.1.1") {
		t.Fatal("expected host to be present")
	}
}

func TestAdd_EmptyHostReturnsError(t *testing.T) {
	w := watchlist.New()
	if err := w.Add("", ""); err == nil {
		t.Fatal("expected error for empty host")
	}
}

func TestAdd_DuplicateIsNoop(t *testing.T) {
	w := watchlist.New()
	_ = w.Add("10.0.0.1", "a")
	_ = w.Add("10.0.0.1", "b")
	if w.Len() != 1 {
		t.Fatalf("expected 1 entry, got %d", w.Len())
	}
}

func TestRemove_ExistingHost(t *testing.T) {
	w := watchlist.New()
	_ = w.Add("10.0.0.1", "")
	if !w.Remove("10.0.0.1") {
		t.Fatal("expected Remove to return true")
	}
	if w.Has("10.0.0.1") {
		t.Fatal("host should no longer be present")
	}
}

func TestRemove_UnknownHostReturnsFalse(t *testing.T) {
	w := watchlist.New()
	if w.Remove("unknown") {
		t.Fatal("expected false for unknown host")
	}
}

func TestAll_PreservesInsertionOrder(t *testing.T) {
	w := watchlist.New()
	hosts := []string{"a.example.com", "b.example.com", "c.example.com"}
	for _, h := range hosts {
		_ = w.Add(h, "")
	}
	entries := w.All()
	if len(entries) != len(hosts) {
		t.Fatalf("expected %d entries, got %d", len(hosts), len(entries))
	}
	for i, e := range entries {
		if e.Host != hosts[i] {
			t.Errorf("position %d: want %s, got %s", i, hosts[i], e.Host)
		}
	}
}

func TestRemove_ReindexesRemainingHosts(t *testing.T) {
	w := watchlist.New()
	_ = w.Add("first", "")
	_ = w.Add("second", "")
	_ = w.Add("third", "")
	w.Remove("second")
	if !w.Has("first") || !w.Has("third") {
		t.Fatal("remaining hosts should still be present")
	}
	if w.Has("second") {
		t.Fatal("removed host should not be present")
	}
	if w.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", w.Len())
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	w := watchlist.New()
	_ = w.Add("x.example.com", "x")
	entries := w.All()
	entries[0].Alias = "mutated"
	if w.All()[0].Alias == "mutated" {
		t.Fatal("All() should return a copy, not a reference")
	}
}
