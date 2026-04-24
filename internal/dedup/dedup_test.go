package dedup_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/dedup"
	"github.com/yourorg/portwatch/internal/scanner"
)

func makeResults(ports ...int) []scanner.Result {
	out := make([]scanner.Result, len(ports))
	for i, p := range ports {
		out[i] = scanner.Result{Host: "localhost", Port: p, Open: true}
	}
	return out
}

func TestChanged_FirstCallIsAlwaysTrue(t *testing.T) {
	s := dedup.New()
	if !s.Changed("host1", makeResults(80, 443)) {
		t.Fatal("expected true on first call")
	}
}

func TestChanged_SameResultsReturnsFalse(t *testing.T) {
	s := dedup.New()
	res := makeResults(80, 443)
	s.Changed("host1", res)
	if s.Changed("host1", res) {
		t.Fatal("expected false for identical consecutive results")
	}
}

func TestChanged_DifferentPortsReturnsTrue(t *testing.T) {
	s := dedup.New()
	s.Changed("host1", makeResults(80))
	if !s.Changed("host1", makeResults(80, 443)) {
		t.Fatal("expected true when port set changes")
	}
}

func TestChanged_DifferentHostsAreIndependent(t *testing.T) {
	s := dedup.New()
	res := makeResults(22)
	s.Changed("hostA", res)
	if !s.Changed("hostB", res) {
		t.Fatal("expected true: hostB has not been seen yet")
	}
}

func TestReset_CausesNextCallToReturnTrue(t *testing.T) {
	s := dedup.New()
	res := makeResults(8080)
	s.Changed("host1", res)
	s.Reset("host1")
	if !s.Changed("host1", res) {
		t.Fatal("expected true after Reset")
	}
}

func TestChanged_EmptyResultsReturnsFalse(t *testing.T) {
	s := dedup.New()
	if s.Changed("host1", nil) {
		t.Fatal("expected false for empty results")
	}
}
