package scheduler_test

import (
	"context"
	"errors"
	"log"
	"io"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scheduler"
)

func silentLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

func TestRun_ExecutesJobImmediately(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	s := scheduler.New(10*time.Second, job, silentLogger())
	s.Run(ctx)

	if atomic.LoadInt32(&count) < 1 {
		t.Error("expected job to be called at least once immediately")
	}
}

func TestRun_ExecutesJobOnTick(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	s := scheduler.New(100*time.Millisecond, job, silentLogger())
	s.Run(ctx)

	if got := atomic.LoadInt32(&count); got < 3 {
		t.Errorf("expected at least 3 executions, got %d", got)
	}
}

func TestRun_ContinuesOnJobError(t *testing.T) {
	var count int32
	job := func(ctx context.Context) error {
		atomic.AddInt32(&count, 1)
		return errors.New("scan failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	s := scheduler.New(100*time.Millisecond, job, silentLogger())
	s.Run(ctx)

	if got := atomic.LoadInt32(&count); got < 2 {
		t.Errorf("expected job to keep running after error, got %d calls", got)
	}
}

func TestNew_NilLogger(t *testing.T) {
	job := func(ctx context.Context) error { return nil }
	s := scheduler.New(time.Second, job, nil)
	if s == nil {
		t.Error("expected non-nil scheduler with nil logger")
	}
}
