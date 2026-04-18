package baseline

import (
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makeResults(ports ...int) []scanner.Result {
	var out []scanner.Result
	for _, p := range ports {
		out = append(out, scanner.Result{Host: "localhost", Port: p, Open: true, Timestamp: time.Now()})
	}
	return out
}

func TestNew_CreatesDir(t *testing.T) {
	dir := t.TempDir() + "/baselines"
	_, err := New(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("dir not created: %v", err)
	}
}

func TestSaveAndLoad(t *testing.T) {
	m, _ := New(t.TempDir())
	results := makeResults(22, 80, 443)

	if err := m.Save("localhost", results); err != nil {
		t.Fatalf("save: %v", err)
	}

	b, err := m.Load("localhost")
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if b == nil {
		t.Fatal("expected baseline, got nil")
	}
	if len(b.Ports) != 3 {
		t.Errorf("expected 3 ports, got %d", len(b.Ports))
	}
}

func TestLoad_NonExistent(t *testing.T) {
	m, _ := New(t.TempDir())
	b, err := m.Load("ghost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b != nil {
		t.Fatal("expected nil for missing baseline")
	}
}

func TestDelete_RemovesFile(t *testing.T) {
	m, _ := New(t.TempDir())
	_ = m.Save("host1", makeResults(80))
	if err := m.Delete("host1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	b, _ := m.Load("host1")
	if b != nil {
		t.Fatal("expected nil after delete")
	}
}

func TestDelete_NonExistent(t *testing.T) {
	m, _ := New(t.TempDir())
	if err := m.Delete("nobody"); err != nil {
		t.Fatalf("delete non-existent should not error: %v", err)
	}
}

func TestSanitize_ReplacesColons(t *testing.T) {
	result := sanitize("192.168.1.1:8080")
	for _, c := range result {
		if c == ':' {
			t.Error("colon not sanitized")
		}
	}
}
