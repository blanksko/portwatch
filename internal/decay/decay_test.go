package decay

import (
	"testing"
	"time"
)

var epoch = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestRecord_IncreasesScore(t *testing.T) {
	tr := New(0.5, time.Minute)
	tr.now = fixedNow(epoch)

	tr.Record("host1", 1.0)
	if got := tr.Score("host1"); got != 1.0 {
		t.Fatalf("want 1.0, got %f", got)
	}
}

func TestScore_UnknownHostIsZero(t *testing.T) {
	tr := New(0.5, time.Minute)
	tr.now = fixedNow(epoch)

	if got := tr.Score("unknown"); got != 0 {
		t.Fatalf("want 0, got %f", got)
	}
}

func TestScore_DecaysAfterInterval(t *testing.T) {
	tr := New(0.5, time.Minute)
	tr.now = fixedNow(epoch)
	tr.Record("host1", 8.0)

	// advance two full intervals → score * 0.5^2 = 2.0
	tr.now = fixedNow(epoch.Add(2 * time.Minute))
	got := tr.Score("host1")
	if got >= 3.0 || got <= 1.0 {
		t.Fatalf("want ~2.0 after two intervals, got %f", got)
	}
}

func TestRecord_AccumulatesMultipleDeltas(t *testing.T) {
	tr := New(0.5, time.Minute)
	tr.now = fixedNow(epoch)

	tr.Record("h", 3.0)
	tr.Record("h", 2.0)
	if got := tr.Score("h"); got != 5.0 {
		t.Fatalf("want 5.0, got %f", got)
	}
}

func TestReset_ClearsScore(t *testing.T) {
	tr := New(0.5, time.Minute)
	tr.now = fixedNow(epoch)

	tr.Record("host1", 10.0)
	tr.Reset("host1")
	if got := tr.Score("host1"); got != 0 {
		t.Fatalf("want 0 after reset, got %f", got)
	}
}

func TestNew_ClampsInvalidDecayRate(t *testing.T) {
	tr := New(-1.0, time.Minute)
	if tr.rate != DefaultDecayRate {
		t.Fatalf("expected rate clamped to default, got %f", tr.rate)
	}

	tr2 := New(0, time.Minute)
	if tr2.rate != DefaultDecayRate {
		t.Fatalf("expected rate clamped to default for zero, got %f", tr2.rate)
	}
}

func TestDifferentHosts_AreIndependent(t *testing.T) {
	tr := New(0.5, time.Minute)
	tr.now = fixedNow(epoch)

	tr.Record("a", 4.0)
	tr.Record("b", 1.0)

	if got := tr.Score("a"); got != 4.0 {
		t.Fatalf("host a: want 4.0, got %f", got)
	}
	if got := tr.Score("b"); got != 1.0 {
		t.Fatalf("host b: want 1.0, got %f", got)
	}
}
