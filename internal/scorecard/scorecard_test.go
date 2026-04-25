package scorecard_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/scorecard"
)

func fixedNow() func() time.Time {
	t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return func() time.Time { return t }
}

func TestGet_UnknownHostIsZero(t *testing.T) {
	sc := scorecard.New()
	s := sc.Get("unknown.host")
	if s.Value != 0 {
		t.Fatalf("expected 0, got %f", s.Value)
	}
}

func TestRecordSuccess_IncreasesScore(t *testing.T) {
	sc := scorecard.New(scorecard.WithGain(10))
	sc.RecordSuccess("host-a")
	if got := sc.Get("host-a").Value; got != 10 {
		t.Fatalf("expected 10, got %f", got)
	}
}

func TestRecordSuccess_CapsAtMax(t *testing.T) {
	sc := scorecard.New(scorecard.WithMax(15), scorecard.WithGain(10))
	sc.RecordSuccess("host-a")
	sc.RecordSuccess("host-a")
	if got := sc.Get("host-a").Value; got != 15 {
		t.Fatalf("expected cap at 15, got %f", got)
	}
}

func TestRecordFailure_DecreasesScore(t *testing.T) {
	sc := scorecard.New(scorecard.WithGain(20), scorecard.WithDecay(5))
	sc.RecordSuccess("host-b")
	sc.RecordFailure("host-b")
	if got := sc.Get("host-b").Value; got != 15 {
		t.Fatalf("expected 15, got %f", got)
	}
}

func TestRecordFailure_FloorsAtZero(t *testing.T) {
	sc := scorecard.New(scorecard.WithDecay(50))
	sc.RecordFailure("host-c")
	if got := sc.Get("host-c").Value; got != 0 {
		t.Fatalf("expected floor at 0, got %f", got)
	}
}

func TestReset_ClearsEntry(t *testing.T) {
	sc := scorecard.New(scorecard.WithGain(10))
	sc.RecordSuccess("host-d")
	sc.Reset("host-d")
	if got := sc.Get("host-d").Value; got != 0 {
		t.Fatalf("expected 0 after reset, got %f", got)
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	sc := scorecard.New(scorecard.WithGain(5))
	sc.RecordSuccess("h1")
	sc.RecordSuccess("h2")
	all := sc.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestAll_MutatingReturnDoesNotAffectInternal(t *testing.T) {
	sc := scorecard.New(scorecard.WithGain(5))
	sc.RecordSuccess("h1")
	all := sc.All()
	all["h1"] = scorecard.Score{Value: 9999}
	if got := sc.Get("h1").Value; got == 9999 {
		t.Fatal("internal state was mutated through returned map")
	}
}

func TestDifferentHosts_AreIndependent(t *testing.T) {
	sc := scorecard.New(scorecard.WithGain(10), scorecard.WithDecay(5))
	sc.RecordSuccess("host-x")
	sc.RecordFailure("host-y")
	if sc.Get("host-x").Value != 10 {
		t.Fatal("host-x score incorrect")
	}
	if sc.Get("host-y").Value != 0 {
		t.Fatal("host-y score should be 0")
	}
}

func TestUpdatedAt_IsSet(t *testing.T) {
	now := time.Now()
	sc := scorecard.New()
	sc.RecordSuccess("ts-host")
	s := sc.Get("ts-host")
	if s.UpdatedAt.Before(now) {
		t.Fatal("UpdatedAt should be set to current time")
	}
}
