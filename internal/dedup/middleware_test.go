package dedup_test

import (
	"context"
	"errors"
	"log/slog"
	"io"
	"testing"

	"github.com/yourorg/portwatch/internal/dedup"
	"github.com/yourorg/portwatch/internal/scanner"
)

func silentLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestGuard_CallsNextOnFirstScan(t *testing.T) {
	called := false
	next := func(_ context.Context, _ string, _ []scanner.Result) error {
		called = true
		return nil
	}
	g := dedup.NewGuard(dedup.New(), next, silentLogger())
	if err := g.Handle(context.Background(), "h1", makeResults(80)); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("expected next to be called")
	}
}

func TestGuard_SkipsNextWhenUnchanged(t *testing.T) {
	calls := 0
	next := func(_ context.Context, _ string, _ []scanner.Result) error {
		calls++
		return nil
	}
	res := makeResults(443)
	g := dedup.NewGuard(dedup.New(), next, silentLogger())
	_ = g.Handle(context.Background(), "h1", res)
	_ = g.Handle(context.Background(), "h1", res)
	if calls != 1 {
		t.Fatalf("expected next called once, got %d", calls)
	}
}

func TestGuard_CallsNextWhenChanged(t *testing.T) {
	calls := 0
	next := func(_ context.Context, _ string, _ []scanner.Result) error {
		calls++
		return nil
	}
	g := dedup.NewGuard(dedup.New(), next, silentLogger())
	_ = g.Handle(context.Background(), "h1", makeResults(80))
	_ = g.Handle(context.Background(), "h1", makeResults(80, 443))
	if calls != 2 {
		t.Fatalf("expected next called twice, got %d", calls)
	}
}

func TestGuard_PropagatesNextError(t *testing.T) {
	want := errors.New("downstream error")
	next := func(_ context.Context, _ string, _ []scanner.Result) error { return want }
	g := dedup.NewGuard(dedup.New(), next, silentLogger())
	if err := g.Handle(context.Background(), "h1", makeResults(22)); !errors.Is(err, want) {
		t.Fatalf("expected %v, got %v", want, err)
	}
}

func TestNewGuard_NilStorePanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil store")
		}
	}()
	dedup.NewGuard(nil, func(_ context.Context, _ string, _ []scanner.Result) error { return nil }, silentLogger())
}
