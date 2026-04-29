package fence_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/fence"
	"github.com/user/portwatch/internal/scanner"
)

func makeResult(port int, open bool) scanner.Result {
	return scanner.Result{
		Host:      "127.0.0.1",
		Port:      port,
		Open:      open,
		Timestamp: time.Now(),
	}
}

func TestNew_ValidRanges(t *testing.T) {
	_, err := fence.New([]fence.Range{{Low: 1, High: 1024}, {Low: 8000, High: 9000}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_InvalidRange_LowExceedsHigh(t *testing.T) {
	_, err := fence.New([]fence.Range{{Low: 1024, High: 80}})
	if err == nil {
		t.Fatal("expected error for inverted range")
	}
}

func TestNew_OutOfBounds(t *testing.T) {
	_, err := fence.New([]fence.Range{{Low: 0, High: 100}})
	if err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestWithin_InsideRange(t *testing.T) {
	f, _ := fence.New([]fence.Range{{Low: 80, High: 443}})
	if !f.Within(80) || !f.Within(443) || !f.Within(200) {
		t.Fatal("expected ports 80, 200, 443 to be within fence")
	}
}

func TestWithin_OutsideRange(t *testing.T) {
	f, _ := fence.New([]fence.Range{{Low: 80, High: 443}})
	if f.Within(79) || f.Within(444) {
		t.Fatal("expected ports 79 and 444 to be outside fence")
	}
}

func TestApply_FiltersToRange(t *testing.T) {
	f, _ := fence.New([]fence.Range{{Low: 80, High: 443}})
	input := []scanner.Result{
		makeResult(22, true),
		makeResult(80, true),
		makeResult(443, true),
		makeResult(8080, true),
	}
	out := f.Apply(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestApply_ClosedPortsDropped(t *testing.T) {
	f, _ := fence.New([]fence.Range{{Low: 1, High: 65535}})
	input := []scanner.Result{
		makeResult(80, false),
		makeResult(443, true),
	}
	out := f.Apply(input)
	if len(out) != 1 || out[0].Port != 443 {
		t.Fatalf("expected only port 443, got %v", out)
	}
}

func TestAdd_ExtendsFence(t *testing.T) {
	f, _ := fence.New([]fence.Range{{Low: 80, High: 80}})
	if err := f.Add(fence.Range{Low: 443, High: 443}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Within(443) {
		t.Fatal("expected 443 to be within fence after Add")
	}
}

func TestRanges_ReturnsCopy(t *testing.T) {
	f, _ := fence.New([]fence.Range{{Low: 1, High: 100}})
	r := f.Ranges()
	r[0].Low = 9999
	if f.Ranges()[0].Low == 9999 {
		t.Fatal("Ranges must return a copy, not a reference")
	}
}
