package sampler_test

import (
	"context"
	"errors"
	"io"
	"log"
	"testing"

	"github.com/user/portwatch/internal/sampler"
	"github.com/user/portwatch/internal/scanner"
)

func silentLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

func TestGuard_SkipsNextWhenBelowThreshold(t *testing.T) {
	s := sampler.New(3)
	called := false
	next := func(_ context.Context, _ string, _ []scanner.Result) error {
		called = true
		return nil
	}
	g := sampler.NewGuard(s, next, silentLogger())

	_ = g.Handle(context.Background(), "host", makeResults("host", 80))
	if called {
		t.Error("next should not be called when threshold not met")
	}
}

func TestGuard_CallsNextWhenThresholdMet(t *testing.T) {
	s := sampler.New(2)
	var received []scanner.Result
	next := func(_ context.Context, _ string, res []scanner.Result) error {
		received = res
		return nil
	}
	g := sampler.NewGuard(s, next, silentLogger())

	for i := 0; i < 2; i++ {
		_ = g.Handle(context.Background(), "host", makeResults("host", 443))
	}
	if len(received) != 1 {
		t.Fatalf("expected 1 stable result, got %d", len(received))
	}
	if received[0].Port != 443 {
		t.Errorf("expected port 443, got %d", received[0].Port)
	}
}

func TestGuard_PropagatesNextError(t *testing.T) {
	s := sampler.New(1)
	sentinel := errors.New("downstream error")
	next := func(_ context.Context, _ string, _ []scanner.Result) error {
		return sentinel
	}
	g := sampler.NewGuard(s, next, silentLogger())

	err := g.Handle(context.Background(), "host", makeResults("host", 22))
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
}

func TestNewGuard_NilSamplerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil sampler")
		}
	}()
	sampler.NewGuard(nil, func(_ context.Context, _ string, _ []scanner.Result) error { return nil }, silentLogger())
}
