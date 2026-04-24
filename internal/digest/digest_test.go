package digest_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/digest"
	"github.com/user/portwatch/internal/scanner"
)

func makeResults(ports ...int) []scanner.Result {
	out := make([]scanner.Result, len(ports))
	for i, p := range ports {
		out[i] = scanner.Result{Host: "host", Port: p, Open: true, Timestamp: time.Now()}
	}
	return out
}

func TestCompute_DeterministicOrder(t *testing.T) {
	a := digest.Compute(makeResults(80, 443, 22))
	b := digest.Compute(makeResults(443, 22, 80))
	if a != b {
		t.Fatalf("expected identical fingerprints, got %s vs %s", a, b)
	}
}

func TestCompute_DifferentPorts(t *testing.T) {
	a := digest.Compute(makeResults(80))
	b := digest.Compute(makeResults(443))
	if a == b {
		t.Fatal("expected different fingerprints for different port sets")
	}
}

func TestCompute_EmptyResults(t *testing.T) {
	a := digest.Compute(nil)
	b := digest.Compute([]scanner.Result{})
	if a != b {
		t.Fatalf("nil and empty should produce the same fingerprint, got %s vs %s", a, b)
	}
}

func TestChanged_FirstCallIsAlwaysChanged(t *testing.T) {
	d := digest.New()
	if !d.Changed("localhost", makeResults(80)) {
		t.Fatal("first call should always report changed")
	}
}

func TestChanged_SameResultsNotChanged(t *testing.T) {
	d := digest.New()
	d.Changed("localhost", makeResults(80, 443))
	if d.Changed("localhost", makeResults(80, 443)) {
		t.Fatal("identical results should not report changed")
	}
}

func TestChanged_DifferentResultsChanged(t *testing.T) {
	d := digest.New()
	d.Changed("localhost", makeResults(80))
	if !d.Changed("localhost", makeResults(80, 443)) {
		t.Fatal("different results should report changed")
	}
}

func TestGet_UnknownHostReturnsEmpty(t *testing.T) {
	d := digest.New()
	if got := d.Get("unknown"); got != "" {
		t.Fatalf("expected empty string, got %q", got)
	}
}

func TestGet_KnownHostReturnsFingerprint(t *testing.T) {
	d := digest.New()
	d.Changed("host", makeResults(22))
	if d.Get("host") == "" {
		t.Fatal("expected non-empty fingerprint after Changed call")
	}
}

func TestReset_CausesNextChangedToReturnTrue(t *testing.T) {
	d := digest.New()
	d.Changed("host", makeResults(80))
	d.Reset("host")
	if !d.Changed("host", makeResults(80)) {
		t.Fatal("after Reset, Changed should return true")
	}
}
