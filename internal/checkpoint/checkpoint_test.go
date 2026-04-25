package checkpoint_test

import (
	"os"
	"testing"
	"time"

	"github.com/example/portwatch/internal/checkpoint"
)

func TestNew_CreatesDir(t *testing.T) {
	dir := t.TempDir()
	subdir := dir + "/nested/checkpoints"
	_, err := checkpoint.New(subdir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := os.Stat(subdir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestSetAndGet(t *testing.T) {
	s := newStore(t)
	now := time.Now().Truncate(time.Second)

	if err := s.Set("host-a", now); err != nil {
		t.Fatalf("Set: %v", err)
	}

	got, ok := s.Get("host-a")
	if !ok {
		t.Fatal("expected checkpoint to exist")
	}
	if !got.Equal(now) {
		t.Fatalf("got %v, want %v", got, now)
	}
}

func TestGet_UnknownHost(t *testing.T) {
	s := newStore(t)
	_, ok := s.Get("unknown")
	if ok {
		t.Fatal("expected no checkpoint for unknown host")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := newStore(t)
	_ = s.Set("host-b", time.Now())

	if err := s.Delete("host-b"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, ok := s.Get("host-b")
	if ok {
		t.Fatal("expected checkpoint to be removed")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := newStore(t)
	_ = s.Set("h1", time.Now())
	_ = s.Set("h2", time.Now())

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestPersistence_ReloadsOnNew(t *testing.T) {
	dir := t.TempDir()
	now := time.Now().Truncate(time.Second).UTC()

	s1, _ := checkpoint.New(dir)
	_ = s1.Set("persist-host", now)

	s2, err := checkpoint.New(dir)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	got, ok := s2.Get("persist-host")
	if !ok {
		t.Fatal("expected persisted checkpoint to be loaded")
	}
	if !got.Equal(now) {
		t.Fatalf("got %v, want %v", got, now)
	}
}

func newStore(t *testing.T) *checkpoint.Store {
	t.Helper()
	s, err := checkpoint.New(t.TempDir())
	if err != nil {
		t.Fatalf("checkpoint.New: %v", err)
	}
	return s
}
