package envelope_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/envelope"
	"github.com/user/portwatch/internal/scanner"
)

func makeResults(ports ...int) []scanner.Result {
	results := make([]scanner.Result, 0, len(ports))
	for _, p := range ports {
		results = append(results, scanner.Result{
			Host: "127.0.0.1",
			Port: p,
			Open: true,
		})
	}
	return results
}

func TestNew_SetsHostAndTime(t *testing.T) {
	before := time.Now().UTC()
	e := envelope.New("10.0.0.1", makeResults(80, 443))
	after := time.Now().UTC()

	if e.Host != "10.0.0.1" {
		t.Errorf("host: got %q, want %q", e.Host, "10.0.0.1")
	}
	if e.ScannedAt.Before(before) || e.ScannedAt.After(after) {
		t.Errorf("ScannedAt %v outside expected range", e.ScannedAt)
	}
	if len(e.Results) != 2 {
		t.Errorf("results: got %d, want 2", len(e.Results))
	}
}

func TestNew_LabelsInitialised(t *testing.T) {
	e := envelope.New("host", nil)
	if e.Labels == nil {
		t.Fatal("Labels map should not be nil after New")
	}
}

func TestWithTags_AppendsAll(t *testing.T) {
	e := envelope.New("host", nil).WithTags("prod", "web")
	if len(e.Tags) != 2 {
		t.Errorf("tags: got %d, want 2", len(e.Tags))
	}
	if e.Tags[0] != "prod" || e.Tags[1] != "web" {
		t.Errorf("unexpected tags: %v", e.Tags)
	}
}

func TestWithLabel_StoresKeyValue(t *testing.T) {
	e := envelope.New("host", nil).WithLabel("env", "staging")
	if got := e.Labels["env"]; got != "staging" {
		t.Errorf("label env: got %q, want %q", got, "staging")
	}
}

func TestWithScanID_SetsField(t *testing.T) {
	e := envelope.New("host", nil).WithScanID("abc-123")
	if e.ScanID != "abc-123" {
		t.Errorf("ScanID: got %q, want %q", e.ScanID, "abc-123")
	}
}

func TestOpenCount_CountsOpenPorts(t *testing.T) {
	results := []scanner.Result{
		{Port: 22, Open: true},
		{Port: 23, Open: false},
		{Port: 80, Open: true},
	}
	e := envelope.New("host", results)
	if got := e.OpenCount(); got != 2 {
		t.Errorf("OpenCount: got %d, want 2", got)
	}
}

func TestOpenCount_Empty(t *testing.T) {
	e := envelope.New("host", nil)
	if got := e.OpenCount(); got != 0 {
		t.Errorf("OpenCount on empty: got %d, want 0", got)
	}
}

func TestChaining_ReturnsSamePointer(t *testing.T) {
	e := envelope.New("host", nil)
	got := e.WithTags("a").WithLabel("k", "v").WithScanID("x")
	if got != e {
		t.Error("chained methods should return the same *Envelope")
	}
}
