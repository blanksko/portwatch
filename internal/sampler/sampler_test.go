package sampler_test

import (
	"testing"

	"github.com/user/portwatch/internal/sampler"
	"github.com/user/portwatch/internal/scanner"
)

func makeResults(host string, ports ...int) []scanner.Result {
	var res []scanner.Result
	for _, p := range ports {
		res = append(res, scanner.Result{Host: host, Port: p, Open: true})
	}
	return res
}

func TestNew_DefaultThresh(t *testing.T) {
	s := sampler.New(0) // zero should clamp to 1
	s.Record(makeResults("localhost", 80))
	got := s.Sample("localhost")
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}

func TestSample_BelowThreshold(t *testing.T) {
	s := sampler.New(3)
	s.Record(makeResults("host1", 443))
	s.Record(makeResults("host1", 443))
	got := s.Sample("host1")
	if len(got) != 0 {
		t.Fatalf("expected 0 results below threshold, got %d", len(got))
	}
}

func TestSample_AtThreshold(t *testing.T) {
	s := sampler.New(3)
	for i := 0; i < 3; i++ {
		s.Record(makeResults("host1", 443))
	}
	got := s.Sample("host1")
	if len(got) != 1 {
		t.Fatalf("expected 1 result at threshold, got %d", len(got))
	}
	if got[0].Port != 443 {
		t.Errorf("expected port 443, got %d", got[0].Port)
	}
}

func TestSample_ClosedPortsIgnored(t *testing.T) {
	s := sampler.New(1)
	closed := []scanner.Result{{Host: "host2", Port: 22, Open: false}}
	s.Record(closed)
	got := s.Sample("host2")
	if len(got) != 0 {
		t.Fatalf("closed ports must not appear in sample, got %d", len(got))
	}
}

func TestSample_UnknownHost(t *testing.T) {
	s := sampler.New(1)
	got := s.Sample("unknown")
	if got != nil {
		t.Errorf("expected nil for unknown host, got %v", got)
	}
}

func TestReset_ClearsHost(t *testing.T) {
	s := sampler.New(1)
	s.Record(makeResults("host3", 8080))
	s.Reset("host3")
	got := s.Sample("host3")
	if len(got) != 0 {
		t.Fatalf("expected empty sample after reset, got %d", len(got))
	}
}

func TestRecord_MultipleHosts(t *testing.T) {
	s := sampler.New(2)
	s.Record(makeResults("a", 80))
	s.Record(makeResults("a", 80))
	s.Record(makeResults("b", 80))

	if len(s.Sample("a")) != 1 {
		t.Error("host a should have 1 sampled port")
	}
	if len(s.Sample("b")) != 0 {
		t.Error("host b should have 0 sampled ports (below threshold)")
	}
}
