package circuit

import (
	"testing"
	"time"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestAllow_NewHost(t *testing.T) {
	b := New(3, time.Minute)
	if !b.Allow("host-a") {
		t.Fatal("expected new host to be allowed")
	}
}

func TestAllow_BelowThreshold(t *testing.T) {
	b := New(3, time.Minute)
	b.RecordFailure("host-a")
	b.RecordFailure("host-a")
	if !b.Allow("host-a") {
		t.Fatal("expected host below threshold to be allowed")
	}
}

func TestAllow_AtThreshold_Blocked(t *testing.T) {
	b := New(3, time.Minute)
	b.RecordFailure("host-a")
	b.RecordFailure("host-a")
	b.RecordFailure("host-a")
	if b.Allow("host-a") {
		t.Fatal("expected host at threshold to be blocked")
	}
}

func TestAllow_AfterCooldown_HalfOpen(t *testing.T) {
	now := time.Now()
	b := New(2, time.Minute)
	b.now = fixedNow(now)
	b.RecordFailure("host-a")
	b.RecordFailure("host-a")
	// advance past cooldown
	b.now = fixedNow(now.Add(2 * time.Minute))
	if !b.Allow("host-a") {
		t.Fatal("expected half-open after cooldown")
	}
	if b.StateOf("host-a") != StateHalfOpen {
		t.Fatalf("expected half-open state, got %s", b.StateOf("host-a"))
	}
}

func TestRecordSuccess_ClosesCiruit(t *testing.T) {
	b := New(2, time.Minute)
	b.RecordFailure("host-a")
	b.RecordFailure("host-a")
	b.RecordSuccess("host-a")
	if b.StateOf("host-a") != StateClosed {
		t.Fatalf("expected closed after success, got %s", b.StateOf("host-a"))
	}
	if !b.Allow("host-a") {
		t.Fatal("expected host to be allowed after success")
	}
}

func TestDifferentHosts_Independent(t *testing.T) {
	b := New(2, time.Minute)
	b.RecordFailure("host-a")
	b.RecordFailure("host-a")
	if !b.Allow("host-b") {
		t.Fatal("expected independent host to be unaffected")
	}
}

func TestStateOf_String(t *testing.T) {
	cases := []struct {
		s    State
		want string
	}{
		{StateClosed, "closed"},
		{StateOpen, "open"},
		{StateHalfOpen, "half-open"},
		{State(99), "unknown"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("State(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}
