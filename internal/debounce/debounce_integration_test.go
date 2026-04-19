package debounce

import (
	"sync"
	"testing"
	"time"
)

func TestDebounce_ConcurrentKeys(t *testing.T) {
	d := New(50 * time.Millisecond)
	var wg sync.WaitGroup
	allowed := make([]bool, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			allowed[idx] = d.Allow("shared-key")
		}(i)
	}
	wg.Wait()
	count := 0
	for _, v := range allowed {
		if v {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 allowed, got %d", count)
	}
}

func TestDebounce_RealTime(t *testing.T) {
	d := New(30 * time.Millisecond)
	if !d.Allow("k") {
		t.Fatal("first call should be allowed")
	}
	if d.Allow("k") {
		t.Fatal("immediate second call should be blocked")
	}
	time.Sleep(40 * time.Millisecond)
	if !d.Allow("k") {
		t.Fatal("call after window should be allowed")
	}
}
