package limiter_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/limiter"
)

func TestNew_DefaultConcurrency(t *testing.T) {
	l := limiter.New(0)
	if l == nil {
		t.Fatal("expected non-nil Limiter")
	}
}

func TestActive_InitiallyZero(t *testing.T) {
	l := limiter.New(5)
	if got := l.Active("host1"); got != 0 {
		t.Fatalf("expected 0 active, got %d", got)
	}
}

func TestAcquireRelease_SingleSlot(t *testing.T) {
	l := limiter.New(1)
	l.Acquire("host1")
	if got := l.Active("host1"); got != 1 {
		t.Fatalf("expected 1 active after Acquire, got %d", got)
	}
	l.Release("host1")
	if got := l.Active("host1"); got != 0 {
		t.Fatalf("expected 0 active after Release, got %d", got)
	}
}

func TestAcquire_BlocksAtLimit(t *testing.T) {
	l := limiter.New(2)
	l.Acquire("h")
	l.Acquire("h")

	done := make(chan struct{})
	go func() {
		l.Acquire("h") // should block until a slot is freed
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("Acquire should have blocked")
	case <-time.After(50 * time.Millisecond):
	}

	l.Release("h")
	select {
	case <-done:
		// expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Acquire did not unblock after Release")
	}
	l.Release("h")
	l.Release("h")
}

func TestDifferentHosts_Independent(t *testing.T) {
	l := limiter.New(1)
	l.Acquire("a")
	l.Acquire("b") // must not block
	if l.Active("a") != 1 || l.Active("b") != 1 {
		t.Fatal("hosts should track independently")
	}
	l.Release("a")
	l.Release("b")
}

func TestReset_ClearsHost(t *testing.T) {
	l := limiter.New(3)
	l.Acquire("x")
	l.Reset("x")
	if got := l.Active("x"); got != 0 {
		t.Fatalf("expected 0 after Reset, got %d", got)
	}
}

func TestConcurrent_NeverExceedsLimit(t *testing.T) {
	const limit = 3
	const workers = 20
	l := limiter.New(limit)

	var peak int64
	var current int64
	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			l.Acquire("shared")
			v := atomic.AddInt64(&current, 1)
			for {
				old := atomic.LoadInt64(&peak)
				if v <= old || atomic.CompareAndSwapInt64(&peak, old, v) {
					break
				}
			}
			time.Sleep(5 * time.Millisecond)
			atomic.AddInt64(&current, -1)
			l.Release("shared")
		}()
	}
	wg.Wait()

	if peak > limit {
		t.Fatalf("peak concurrency %d exceeded limit %d", peak, limit)
	}
}
