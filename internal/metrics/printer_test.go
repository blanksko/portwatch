package metrics_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestPrint_ContainsExpectedFields(t *testing.T) {
	c := metrics.New()
	c.Record(8, 42*time.Millisecond)
	s := c.Snapshot()

	var buf bytes.Buffer
	if err := metrics.Print(&buf, s); err != nil {
		t.Fatalf("Print returned error: %v", err)
	}

	out := buf.String()
	for _, want := range []string{"Total scans", "Open ports", "Last scan", "Last duration", "1", "8"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n%s", want, out)
		}
	}
}

func TestPrint_NoScans(t *testing.T) {
	c := metrics.New()
	s := c.Snapshot()

	var buf bytes.Buffer
	if err := metrics.Print(&buf, s); err != nil {
		t.Fatalf("Print returned error: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "-") {
		t.Errorf("expected '-' for missing last scan, got:\n%s", out)
	}
}
