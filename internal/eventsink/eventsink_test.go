package eventsink_test

import (
	"errors"
	"log"
	"io"
	"testing"

	"github.com/example/portwatch/internal/eventsink"
	"github.com/example/portwatch/internal/scanner"
	"github.com/example/portwatch/internal/snapshot"
)

func silentLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

func makeDiff(opened, closed []int) snapshot.Diff {
	toResults := func(ports []int) []scanner.Result {
		var out []scanner.Result
		for _, p := range ports {
			out = append(out, scanner.Result{Host: "localhost", Port: p, Open: true})
		}
		return out
	}
	return snapshot.Diff{
		Opened: toResults(opened),
		Closed: toResults(closed),
	}
}

func TestDispatch_NoHandlers(t *testing.T) {
	s := eventsink.New(silentLogger())
	if err := s.Dispatch(makeDiff([]int{80}, nil)); err != nil {
		t.Fatalf("expected no error with no handlers, got %v", err)
	}
}

func TestDispatch_CallsAllHandlers(t *testing.T) {
	s := eventsink.New(silentLogger())
	called := map[string]bool{}
	s.Register("a", func(d snapshot.Diff) error { called["a"] = true; return nil })
	s.Register("b", func(d snapshot.Diff) error { called["b"] = true; return nil })

	if err := s.Dispatch(makeDiff([]int{443}, nil)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, name := range []string{"a", "b"} {
		if !called[name] {
			t.Errorf("handler %q was not called", name)
		}
	}
}

func TestDispatch_HandlerError_ContinuesAndReturnsError(t *testing.T) {
	s := eventsink.New(silentLogger())
	var secondCalled bool
	s.Register("fail", func(d snapshot.Diff) error { return errors.New("boom") })
	s.Register("ok", func(d snapshot.Diff) error { secondCalled = true; return nil })

	err := s.Dispatch(makeDiff(nil, []int{22}))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !secondCalled {
		t.Error("second handler should have been called despite first error")
	}
}

func TestLen_TracksRegistrations(t *testing.T) {
	s := eventsink.New(silentLogger())
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
	s.Register("x", func(d snapshot.Diff) error { return nil })
	s.Register("y", func(d snapshot.Diff) error { return nil })
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
}

func TestNew_NilLoggerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil logger")
		}
	}()
	eventsink.New(nil)
}

func TestRegister_NilHandlerPanics(t *testing.T) {
	s := eventsink.New(silentLogger())
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil handler")
		}
	}()
	s.Register("bad", nil)
}
