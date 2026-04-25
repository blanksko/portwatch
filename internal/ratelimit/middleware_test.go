package ratelimit

import (
	"errors"
	"log"
	"io"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func silentLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

func makeNext(results []scanner.Result, err error) func(string) ([]scanner.Result, error) {
	return func(_ string) ([]scanner.Result, error) {
		return results, err
	}
}

func TestGuard_AllowsFirstCall(t *testing.T) {
	r := New(100 * time.Millisecond)
	results := []scanner.Result{{Host: "localhost", Port: 80, Open: true}}
	g := NewGuard(r, silentLogger(), makeNext(results, nil))

	got, err := g.Run("localhost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}

func TestGuard_BlocksSecondCall(t *testing.T) {
	r := New(100 * time.Millisecond)
	g := NewGuard(r, silentLogger(), makeNext([]scanner.Result{{Host: "h", Port: 22, Open: true}}, nil))

	g.Run("h") // first — allowed
	got, err := g.Run("h") // second — blocked
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Fatalf("expected nil results when rate-limited, got %v", got)
	}
}

func TestGuard_AllowsAfterReset(t *testing.T) {
	r := New(10 * time.Second)
	results := []scanner.Result{{Host: "h", Port: 443, Open: true}}
	g := NewGuard(r, silentLogger(), makeNext(results, nil))

	g.Run("h")
	g.Reset("h")
	got, err := g.Run("h")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected result after reset, got %d", len(got))
	}
}

func TestGuard_PropagatesNextError(t *testing.T) {
	r := New(100 * time.Millisecond)
	expected := errors.New("scan failed")
	g := NewGuard(r, silentLogger(), makeNext(nil, expected))

	_, err := g.Run("host")
	if !errors.Is(err, expected) {
		t.Fatalf("expected %v, got %v", expected, err)
	}
}

func TestNewGuard_NilRateLimitPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil RateLimit")
		}
	}()
	NewGuard(nil, silentLogger(), makeNext(nil, nil))
}

func TestNewGuard_NilLoggerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil logger")
		}
	}()
	NewGuard(New(time.Second), nil, makeNext(nil, nil))
}
