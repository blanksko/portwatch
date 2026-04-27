package prestige_test

import (
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/prestige"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestGet_UnknownHost(t *testing.T) {
	tr := prestige.New()
	_, ok := tr.Get("unknown.host")
	if ok {
		t.Fatal("expected false for unknown host")
	}
}

func TestRecord_IncreasesScore(t *testing.T) {
	tr := prestige.New(prestige.WithGain(2.0), prestige.WithNow(fixedNow(epoch)))
	score := tr.Record("host-a")
	if score != 2.0 {
		t.Fatalf("expected 2.0, got %f", score)
	}
	score = tr.Record("host-a")
	if score != 4.0 {
		t.Fatalf("expected 4.0 after second deviation, got %f", score)
	}
}

func TestRecord_TracksDeviations(t *testing.T) {
	tr := prestige.New(prestige.WithNow(fixedNow(epoch)))
	tr.Record("host-b")
	tr.Record("host-b")
	r, ok := tr.Get("host-b")
	if !ok {
		t.Fatal("expected record")
	}
	if r.Deviations != 2 {
		t.Fatalf("expected 2 deviations, got %d", r.Deviations)
	}
}

func TestDecay_ReducesScore(t *testing.T) {
	tr := prestige.New(prestige.WithGain(10.0), prestige.WithDecay(0.5), prestige.WithNow(fixedNow(epoch)))
	tr.Record("host-c") // score = 10
	tr.Decay()           // score = 10 - 5 = 5
	r, _ := tr.Get("host-c")
	if r.Score != 5.0 {
		t.Fatalf("expected 5.0 after decay, got %f", r.Score)
	}
}

func TestDecay_NeverGoesNegative(t *testing.T) {
	tr := prestige.New(prestige.WithGain(1.0), prestige.WithDecay(2.0), prestige.WithNow(fixedNow(epoch)))
	tr.Record("host-d")
	tr.Decay()
	r, _ := tr.Get("host-d")
	if r.Score < 0 {
		t.Fatalf("score should not be negative, got %f", r.Score)
	}
}

func TestReset_ClearsRecord(t *testing.T) {
	tr := prestige.New(prestige.WithNow(fixedNow(epoch)))
	tr.Record("host-e")
	tr.Reset("host-e")
	_, ok := tr.Get("host-e")
	if ok {
		t.Fatal("expected record to be removed after reset")
	}
}

func TestRecord_DifferentHostsAreIndependent(t *testing.T) {
	tr := prestige.New(prestige.WithGain(3.0), prestige.WithNow(fixedNow(epoch)))
	tr.Record("alpha")
	tr.Record("alpha")
	tr.Record("beta")

	a, _ := tr.Get("alpha")
	b, _ := tr.Get("beta")

	if a.Score != 6.0 {
		t.Fatalf("alpha: expected 6.0, got %f", a.Score)
	}
	if b.Score != 3.0 {
		t.Fatalf("beta: expected 3.0, got %f", b.Score)
	}
}
