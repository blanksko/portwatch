package shadow_test

import (
	"os"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/shadow"
)

func makeResults(ports ...int) []scanner.Result {
	out := make([]scanner.Result, len(ports))
	for i, p := range ports {
		out[i] = scanner.Result{Port: p, Open: true}
	}
	return out
}

func TestNew_CreatesDir(t *testing.T) {
	dir := t.TempDir() + "/shadow"
	_, err := shadow.New(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("directory was not created")
	}
}

func TestRecordAndLoad(t *testing.T) {
	store, _ := shadow.New(t.TempDir())
	results := makeResults(22, 80, 443)

	if err := store.Record("192.168.1.1", results); err != nil {
		t.Fatalf("Record: %v", err)
	}
	rec, err := store.Load("192.168.1.1")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if rec.Host != "192.168.1.1" {
		t.Errorf("host = %q, want %q", rec.Host, "192.168.1.1")
	}
	if len(rec.Results) != 3 {
		t.Errorf("results len = %d, want 3", len(rec.Results))
	}
}

func TestLoad_NonExistent(t *testing.T) {
	store, _ := shadow.New(t.TempDir())
	_, err := store.Load("10.0.0.99")
	if !os.IsNotExist(err) {
		t.Errorf("expected os.ErrNotExist, got %v", err)
	}
}

func TestDelete_RemovesRecord(t *testing.T) {
	store, _ := shadow.New(t.TempDir())
	_ = store.Record("10.0.0.1", makeResults(8080))
	if err := store.Delete("10.0.0.1"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := store.Load("10.0.0.1")
	if !os.IsNotExist(err) {
		t.Errorf("expected ErrNotExist after delete, got %v", err)
	}
}

func TestDelete_UnknownHostIsNoop(t *testing.T) {
	store, _ := shadow.New(t.TempDir())
	if err := store.Delete("ghost"); err != nil {
		t.Errorf("expected no error for unknown host, got %v", err)
	}
}

func TestCompare_NoChanges(t *testing.T) {
	store, _ := shadow.New(t.TempDir())
	results := makeResults(22, 443)
	_ = store.Record("host1", results)
	rec, _ := store.Load("host1")

	delta := shadow.Compare(rec, results)
	if delta.HasChanges() {
		t.Errorf("expected no changes, got added=%v removed=%v", delta.Added, delta.Removed)
	}
}

func TestCompare_AddedPort(t *testing.T) {
	store, _ := shadow.New(t.TempDir())
	_ = store.Record("host2", makeResults(22))
	rec, _ := store.Load("host2")

	delta := shadow.Compare(rec, makeResults(22, 8080))
	if len(delta.Added) != 1 || delta.Added[0] != 8080 {
		t.Errorf("Added = %v, want [8080]", delta.Added)
	}
	if len(delta.Removed) != 0 {
		t.Errorf("Removed = %v, want []", delta.Removed)
	}
}

func TestCompare_RemovedPort(t *testing.T) {
	store, _ := shadow.New(t.TempDir())
	_ = store.Record("host3", makeResults(22, 80))
	rec, _ := store.Load("host3")

	delta := shadow.Compare(rec, makeResults(22))
	if len(delta.Removed) != 1 || delta.Removed[0] != 80 {
		t.Errorf("Removed = %v, want [80]", delta.Removed)
	}
}
