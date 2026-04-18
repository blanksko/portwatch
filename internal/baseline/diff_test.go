package baseline

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func baselineWith(host string, ports ...int) *Baseline {
	var results []scanner.Result
	for _, p := range ports {
		results = append(results, scanner.Result{Host: host, Port: p, Open: true, Timestamp: time.Now()})
	}
	return &Baseline{Host: host, Ports: results}
}

func liveResults(host string, ports ...int) []scanner.Result {
	var out []scanner.Result
	for _, p := range ports {
		out = append(out, scanner.Result{Host: host, Port: p, Open: true, Timestamp: time.Now()})
	}
	return out
}

func TestCompare_NoDeviation(t *testing.T) {
	b := baselineWith("host", 80, 443)
	live := liveResults("host", 80, 443)
	if dev := Compare(b, live); dev != nil {
		t.Errorf("expected nil deviation, got %+v", dev)
	}
}

func TestCompare_ExtraPort(t *testing.T) {
	b := baselineWith("host", 80)
	live := liveResults("host", 80, 8080)
	dev := Compare(b, live)
	if dev == nil {
		t.Fatal("expected deviation")
	}
	if len(dev.Extra) != 1 || dev.Extra[0] != 8080 {
		t.Errorf("unexpected extra ports: %v", dev.Extra)
	}
}

func TestCompare_MissingPort(t *testing.T) {
	b := baselineWith("host", 80, 443)
	live := liveResults("host", 80)
	dev := Compare(b, live)
	if dev == nil {
		t.Fatal("expected deviation")
	}
	if len(dev.Missing) != 1 || dev.Missing[0] != 443 {
		t.Errorf("unexpected missing ports: %v", dev.Missing)
	}
}

func TestCompare_NilBaseline(t *testing.T) {
	live := liveResults("host", 80)
	if dev := Compare(nil, live); dev != nil {
		t.Error("expected nil when no baseline")
	}
}
