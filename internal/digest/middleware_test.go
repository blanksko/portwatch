package digest_test

import (
	"errors"
	"io"
	"log"
	"testing"

	"github.com/user/portwatch/internal/digest"
	"github.com/user/portwatch/internal/scanner"
)

func silentLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

func TestGuard_CallsNextOnFirstScan(t *testing.T) {
	called := false
	g := digest.NewGuard(digest.New(), silentLogger(), func(_ string, _ []scanner.Result) error {
		called = true
		return nil
	})
	if err := g.Handle("host", makeResults(80)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected next to be called on first scan")
	}
}

func TestGuard_SkipsNextWhenUnchanged(t *testing.T) {
	calls := 0
	g := digest.NewGuard(digest.New(), silentLogger(), func(_ string, _ []scanner.Result) error {
		calls++
		return nil
	})
	res := makeResults(80, 443)
	g.Handle("host", res) //nolint:errcheck
	g.Handle("host", res) //nolint:errcheck
	if calls != 1 {
		t.Fatalf("expected next to be called once, got %d", calls)
	}
}

func TestGuard_CallsNextWhenChanged(t *testing.T) {
	calls := 0
	g := digest.NewGuard(digest.New(), silentLogger(), func(_ string, _ []scanner.Result) error {
		calls++
		return nil
	})
	g.Handle("host", makeResults(80))  //nolint:errcheck
	g.Handle("host", makeResults(443)) //nolint:errcheck
	if calls != 2 {
		t.Fatalf("expected next to be called twice, got %d", calls)
	}
}

func TestGuard_PropagatesNextError(t *testing.T) {
	sentinel := errors.New("downstream failure")
	g := digest.NewGuard(digest.New(), silentLogger(), func(_ string, _ []scanner.Result) error {
		return sentinel
	})
	if err := g.Handle("host", makeResults(80)); !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestNewGuard_NilDigestPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil Digest")
		}
	}()
	digest.NewGuard(nil, silentLogger(), func(_ string, _ []scanner.Result) error { return nil })
}
