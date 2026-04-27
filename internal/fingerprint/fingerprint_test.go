package fingerprint_test

import (
	"testing"

	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/scanner"
)

func makeResults(host string, ports []int, open bool) []scanner.Result {
	out := make([]scanner.Result, len(ports))
	for i, p := range ports {
		out[i] = scanner.Result{Host: host, Port: p, Open: open}
	}
	return out
}

func TestCompute_EmptySlice(t *testing.T) {
	f := fingerprint.Compute(nil)
	if f.String() == "" {
		t.Fatal("expected non-empty fingerprint for nil input")
	}
}

func TestCompute_OnlyClosedPorts_MatchesEmpty(t *testing.T) {
	closed := makeResults("host-a", []int{80, 443}, false)
	got := fingerprint.Compute(closed)
	want := fingerprint.Compute(nil)
	if !fingerprint.Equal(got, want) {
		t.Errorf("closed-only scan should equal empty fingerprint; got %s want %s", got, want)
	}
}

func TestCompute_Deterministic(t *testing.T) {
	results := []scanner.Result{
		{Host: "10.0.0.1", Port: 443, Open: true},
		{Host: "10.0.0.1", Port: 80, Open: true},
		{Host: "10.0.0.1", Port: 22, Open: true},
	}
	a := fingerprint.Compute(results)

	// Reverse order — fingerprint must be the same.
	reversed := []scanner.Result{results[2], results[1], results[0]}
	b := fingerprint.Compute(reversed)

	if !fingerprint.Equal(a, b) {
		t.Errorf("fingerprint is not deterministic: %s vs %s", a, b)
	}
}

func TestCompute_DifferentPorts_DifferentFingerprint(t *testing.T) {
	a := fingerprint.Compute(makeResults("host", []int{80}, true))
	b := fingerprint.Compute(makeResults("host", []int{443}, true))
	if fingerprint.Equal(a, b) {
		t.Error("different open ports should produce different fingerprints")
	}
}

func TestCompute_DifferentHosts_DifferentFingerprint(t *testing.T) {
	a := fingerprint.Compute(makeResults("host-a", []int{80}, true))
	b := fingerprint.Compute(makeResults("host-b", []int{80}, true))
	if fingerprint.Equal(a, b) {
		t.Error("same ports on different hosts should produce different fingerprints")
	}
}

func TestCompute_ClosedPortsIgnored(t *testing.T) {
	open := makeResults("host", []int{22, 80}, true)
	mixed := append(open, makeResults("host", []int{8080, 9090}, false)...)

	if !fingerprint.Equal(fingerprint.Compute(open), fingerprint.Compute(mixed)) {
		t.Error("closed ports should not affect the fingerprint")
	}
}

func TestEqual_SameValue(t *testing.T) {
	f := fingerprint.Compute(makeResults("h", []int{22}, true))
	if !fingerprint.Equal(f, f) {
		t.Error("fingerprint should equal itself")
	}
}
