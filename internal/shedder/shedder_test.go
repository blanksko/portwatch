package shedder_test

import (
	"sync"
	"testing"

	"github.com/user/portwatch/internal/shedder"
)

func TestAcquire_BelowLimit(t *testing.T) {
	s := shedder.New(3)
	if err := s.Acquire(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if got := s.Active(); got != 1 {
		t.Fatalf("active = %d, want 1", got)
	}
}

func TestAcquire_AtLimit_Blocked(t *testing.T) {
	s := shedder.New(2)
	_ = s.Acquire()
	_ = s.Acquire()
	if err := s.Acquire(); err != shedder.ErrShed {
		t.Fatalf("expected ErrShed, got %v", err)
	}
}

func TestRelease_DecrementsActive(t *testing.T) {
	s := shedder.New(2)
	_ = s.Acquire()
	_ = s.Acquire()
	s.Release()
	if got := s.Active(); got != 1 {
		t.Fatalf("active = %d, want 1", got)
	}
}

func TestRelease_AfterFullDrain_NoUnderflow(t *testing.T) {
	s := shedder.New(1)
	_ = s.Acquire()
	s.Release()
	s.Release() // extra release should not go negative
	if got := s.Active(); got != 0 {
		t.Fatalf("active = %d, want 0", got)
	}
}

func TestAcquireRelease_AllowsReuse(t *testing.T) {
	s := shedder.New(1)
	_ = s.Acquire()
	s.Release()
	if err := s.Acquire(); err != nil {
		t.Fatalf("expected slot to be free after release, got %v", err)
	}
}

func TestNew_PanicsOnZeroLimit(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for limit=0")
		}
	}()
	shedder.New(0)
}

func TestConcurrent_NeverExceedsLimit(t *testing.T) {
	const limit = 5
	const goroutines = 50
	s := shedder.New(limit)
	var mu sync.Mutex
	var maxSeen int
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			if err := s.Acquire(); err != nil {
				return
			}
			defer s.Release()
			mu.Lock()
			if a := s.Active(); a > maxSeen {
				maxSeen = a
			}
			mu.Unlock()
		}()
	}
	wg.Wait()
	if maxSeen > limit {
		t.Fatalf("active exceeded limit: got %d, limit %d", maxSeen, limit)
	}
}
