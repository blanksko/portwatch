package filter_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
)

func makeResults(ports ...int) []scanner.Result {
	results := make([]scanner.Result, len(ports))
	for i, p := range ports {
		results[i] = scanner.Result{Host: "localhost", Port: p, Open: true, Timestamp: time.Now()}
	}
	return results
}

func TestApply_NoRules_AllowsAll(t *testing.T) {
	f := filter.New(nil, nil)
	in := makeResults(80, 443, 8080)
	out := f.Apply(in)
	if len(out) != 3 {
		t.Fatalf("expected 3 results, got %d", len(out))
	}
}

func TestApply_IncludeList(t *testing.T) {
	f := filter.New([]int{80, 443}, nil)
	in := makeResults(80, 443, 8080)
	out := f.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
}

func TestApply_ExcludeList(t *testing.T) {
	f := filter.New(nil, []int{8080})
	in := makeResults(80, 443, 8080)
	out := f.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out))
	}
	for _, r := range out {
		if r.Port == 8080 {
			t.Fatal("excluded port 8080 present in output")
		}
	}
}

func TestApply_ExcludeTakesPrecedence(t *testing.T) {
	f := filter.New([]int{80, 8080}, []int{8080})
	in := makeResults(80, 8080)
	out := f.Apply(in)
	if len(out) != 1 || out[0].Port != 80 {
		t.Fatalf("expected only port 80, got %v", out)
	}
}

func TestApply_EmptyInput(t *testing.T) {
	f := filter.New([]int{80}, nil)
	out := f.Apply([]scanner.Result{})
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d", len(out))
	}
}
