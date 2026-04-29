package ledger

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeResults(host string, ports []int, open bool) []scanner.Result {
	out := make([]scanner.Result, len(ports))
	for i, p := range ports {
		out[i] = scanner.Result{Host: host, Port: p, Open: open}
	}
	return out
}

func TestRecord_CountsOpenPorts(t *testing.T) {
	l := New()
	l.Record(makeResults("host-a", []int{22, 80, 443}, true))
	e, ok := l.Get("host-a")
	if !ok {
		t.Fatal("expected entry for host-a")
	}
	if e.OpenPorts != 3 {
		t.Fatalf("want 3 open ports, got %d", e.OpenPorts)
	}
}

func TestRecord_IgnoresClosedPorts(t *testing.T) {
	l := New()
	l.Record(makeResults("host-b", []int{8080, 9090}, false))
	_, ok := l.Get("host-b")
	if ok {
		t.Fatal("closed ports should not create an entry")
	}
}

func TestRecord_Accumulates(t *testing.T) {
	l := New()
	l.Record(makeResults("host-c", []int{22}, true))
	l.Record(makeResults("host-c", []int{80}, true))
	e, _ := l.Get("host-c")
	if e.OpenPorts != 2 {
		t.Fatalf("want 2, got %d", e.OpenPorts)
	}
}

func TestGet_UnknownHost(t *testing.T) {
	l := New()
	_, ok := l.Get("unknown")
	if ok {
		t.Fatal("expected no entry for unknown host")
	}
}

func TestReset_ClearsHost(t *testing.T) {
	l := New()
	l.Record(makeResults("host-d", []int{22}, true))
	l.Reset("host-d")
	_, ok := l.Get("host-d")
	if ok {
		t.Fatal("entry should have been removed")
	}
}

func TestReset_ClearsAll(t *testing.T) {
	l := New()
	l.Record(makeResults("host-e", []int{22}, true))
	l.Record(makeResults("host-f", []int{80}, true))
	l.Reset("")
	if entries := l.All(); len(entries) != 0 {
		t.Fatalf("expected empty ledger, got %d entries", len(entries))
	}
}

func TestRecord_SetsLastSeen(t *testing.T) {
	before := time.Now()
	l := New()
	l.Record(makeResults("host-g", []int{443}, true))
	e, _ := l.Get("host-g")
	if e.LastSeen.Before(before) {
		t.Fatal("LastSeen should be at or after test start")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	l := New()
	l.Record(makeResults("h1", []int{22}, true))
	l.Record(makeResults("h2", []int{80}, true))
	if got := len(l.All()); got != 2 {
		t.Fatalf("want 2 entries, got %d", got)
	}
}
