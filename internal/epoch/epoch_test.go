package epoch

import (
	"sync"
	"testing"
)

func TestBump_StartsAtOne(t *testing.T) {
	tr := New()
	got := tr.Bump("host-a")
	if got != 1 {
		t.Fatalf("expected epoch 1, got %d", got)
	}
}

func TestBump_Increments(t *testing.T) {
	tr := New()
	tr.Bump("host-a")
	got := tr.Bump("host-a")
	if got != 2 {
		t.Fatalf("expected epoch 2 after second bump, got %d", got)
	}
}

func TestGet_UnknownHostIsZero(t *testing.T) {
	tr := New()
	if got := tr.Get("unknown"); got != 0 {
		t.Fatalf("expected 0 for unknown host, got %d", got)
	}
}

func TestGet_ReturnsCurrentEpoch(t *testing.T) {
	tr := New()
	tr.Bump("host-b")
	tr.Bump("host-b")
	if got := tr.Get("host-b"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestReset_ClearsEpoch(t *testing.T) {
	tr := New()
	tr.Bump("host-c")
	tr.Reset("host-c")
	if got := tr.Get("host-c"); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestSnapshot_ReturnsCopy(t *testing.T) {
	tr := New()
	tr.Bump("h1")
	tr.Bump("h1")
	tr.Bump("h2")
	snap := tr.Snapshot()
	if snap["h1"] != 2 {
		t.Fatalf("expected h1=2, got %d", snap["h1"])
	}
	if snap["h2"] != 1 {
		t.Fatalf("expected h2=1, got %d", snap["h2"])
	}
	// Mutating the snapshot must not affect the tracker.
	snap["h1"] = 99
	if tr.Get("h1") != 2 {
		t.Fatal("snapshot mutation affected tracker")
	}
}

func TestStale_DetectsOldEpoch(t *testing.T) {
	tr := New()
	ep := tr.Bump("host-d") // ep == 1
	tr.Bump("host-d")       // tracker now at 2
	if !tr.Stale("host-d", ep) {
		t.Fatal("expected epoch 1 to be stale when tracker is at 2")
	}
}

func TestStale_CurrentEpochIsNotStale(t *testing.T) {
	tr := New()
	ep := tr.Bump("host-e")
	if tr.Stale("host-e", ep) {
		t.Fatal("current epoch should not be considered stale")
	}
}

func TestBump_ConcurrentSafety(t *testing.T) {
	tr := New()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tr.Bump("concurrent-host")
		}()
	}
	wg.Wait()
	if got := tr.Get("concurrent-host"); got != 50 {
		t.Fatalf("expected 50 after 50 concurrent bumps, got %d", got)
	}
}
