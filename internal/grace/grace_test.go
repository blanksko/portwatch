package grace_test

import (
	"context"
	"log"
	"io"
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/grace"
)

func silentLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}

func TestNew_NilLoggerPanics(t *testing.T) {
	t.Parallel()
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil logger")
		}
	}()
	grace.New(time.Second, nil)
}

func TestAcquireRelease_TracksActive(t *testing.T) {
	t.Parallel()
	c := grace.New(time.Second, silentLogger())

	if !c.Acquire() {
		t.Fatal("expected Acquire to return true before shutdown")
	}
	if got := c.Active(); got != 1 {
		t.Fatalf("expected 1 active, got %d", got)
	}
	c.Release()
	if got := c.Active(); got != 0 {
		t.Fatalf("expected 0 active after release, got %d", got)
	}
}

func TestAcquire_ReturnsFalseAfterShutdown(t *testing.T) {
	t.Parallel()
	c := grace.New(50*time.Millisecond, silentLogger())

	ctx, cancel := context.WithCancel(context.Background())
	_ = c.Wait(ctx)
	// trigger shutdown via parent cancellation
	cancel()
	time.Sleep(100 * time.Millisecond)

	if c.Acquire() {
		t.Fatal("expected Acquire to return false after shutdown")
	}
}

func TestWait_CancelledByParent(t *testing.T) {
	t.Parallel()
	c := grace.New(200*time.Millisecond, silentLogger())

	ctx, cancel := context.WithCancel(context.Background())
	shutCtx := c.Wait(ctx)
	cancel()

	select {
	case <-shutCtx.Done():
		// expected
	case <-time.After(500 * time.Millisecond):
		t.Fatal("shutdown context was not cancelled within timeout")
	}
}

func TestDrain_WaitsForActiveWork(t *testing.T) {
	t.Parallel()
	c := grace.New(2*time.Second, silentLogger())

	ctx, cancel := context.WithCancel(context.Background())
	shutCtx := c.Wait(ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	if !c.Acquire() {
		t.Fatal("acquire failed before shutdown")
	}

	go func() {
		defer wg.Done()
		time.Sleep(150 * time.Millisecond)
		c.Release()
	}()

	cancel() // trigger shutdown

	// shutCtx should not be done until Release is called
	select {
	case <-shutCtx.Done():
		// ok – drain finished
	case <-time.After(1 * time.Second):
		t.Fatal("shutdown did not complete after work finished")
	}
	wg.Wait()
}
