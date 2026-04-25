package pipeline_test

import (
	"context"
	"errors"
	"log"
	"io"
	"testing"

	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
)

func silentLogger() *log.Logger { return log.New(io.Discard, "", 0) }

func makeResults(ports ...int) []scanner.Result {
	out := make([]scanner.Result, len(ports))
	for i, p := range ports {
		out[i] = scanner.Result{Host: "localhost", Port: p, Open: true}
	}
	return out
}

func TestRun_CallsTerminalWithResults(t *testing.T) {
	p := pipeline.New(silentLogger())
	results := makeResults(80, 443)
	var got []scanner.Result
	err := p.Run(context.Background(), results, func(r []scanner.Result) error {
		got = r
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestRun_StagesExecuteInOrder(t *testing.T) {
	p := pipeline.New(silentLogger())
	var order []string
	for _, name := range []string{"a", "b", "c"} {
		n := name
		p.Use(pipeline.Stage{
			Name: n,
			Handler: func(_ context.Context, r []scanner.Result, next func([]scanner.Result) error) error {
				order = append(order, n)
				return next(r)
			},
		})
	}
	_ = p.Run(context.Background(), nil, func(r []scanner.Result) error { return nil })
	if len(order) != 3 || order[0] != "a" || order[1] != "b" || order[2] != "c" {
		t.Fatalf("unexpected order: %v", order)
	}
}

func TestRun_StageErrorHaltsChain(t *testing.T) {
	p := pipeline.New(silentLogger())
	sentinel := errors.New("boom")
	p.Use(pipeline.Stage{
		Name: "fail",
		Handler: func(_ context.Context, _ []scanner.Result, _ func([]scanner.Result) error) error {
			return sentinel
		},
	})
	called := false
	err := p.Run(context.Background(), nil, func(_ []scanner.Result) error {
		called = true
		return nil
	})
	if called {
		t.Fatal("terminal should not have been called")
	}
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestNew_NilLoggerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil logger")
		}
	}()
	pipeline.New(nil)
}
