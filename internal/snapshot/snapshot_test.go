package snapshot_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func makeResults(ports []int, open bool) []scanner.ScanResult {
	var results []scanner.ScanResult
	for _, p := range ports {
		results = append(results, scanner.ScanResult{
			Port:      p,
			Open:      open,
			Timestamp: time.Now(),
		})
	}
	return results
}

func TestSaveAndLoad(t *testing.T) {
	tmp, err := os.CreateTemp("", "portwatch-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Close()
	defer os.Remove(tmp.Name())

	snap := snapshot.New("localhost", makeResults([]int{80, 443}, true))
	if err := snapshot.Save(tmp.Name(), snap); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := snapshot.Load(tmp.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if loaded.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", loaded.Host)
	}
	if len(loaded.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(loaded.Results))
	}
}

func TestLoad_NonExistent(t *testing.T) {
	snap, err := snapshot.Load("/nonexistent/path/snap.json")
	if err != nil {
		t.Fatalf("expected nil error for missing file, got %v", err)
	}
	if snap != nil {
		t.Error("expected nil snapshot for missing file")
	}
}

func TestDiff_NewPorts(t *testing.T) {
	prev := snapshot.New("localhost", makeResults([]int{80}, true))
	curr := snapshot.New("localhost", makeResults([]int{80, 443}, true))

	changes := snapshot.Diff(prev, curr)
	if len(changes) != 1 || changes[0].Port != 443 || changes[0].Change != snapshot.ChangeOpened {
		t.Errorf("unexpected changes: %+v", changes)
	}
}

func TestDiff_ClosedPort(t *testing.T) {
	prev := snapshot.New("localhost", makeResults([]int{80, 8080}, true))
	curr := snapshot.New("localhost", makeResults([]int{80}, true))

	changes := snapshot.Diff(prev, curr)
	if len(changes) != 1 || changes[0].Port != 8080 || changes[0].Change != snapshot.ChangeClosed {
		t.Errorf("unexpected changes: %+v", changes)
	}
}

func TestDiff_NilPrev(t *testing.T) {
	curr := snapshot.New("localhost", makeResults([]int{22, 80}, true))
	changes := snapshot.Diff(nil, curr)
	if len(changes) != 2 {
		t.Errorf("expected 2 opened changes, got %d", len(changes))
	}
}
