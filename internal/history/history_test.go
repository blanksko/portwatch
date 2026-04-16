package history_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/history"
	"github.com/user/portwatch/internal/scanner"
)

func makeResults() []scanner.Result {
	return []scanner.Result{
		{Host: "localhost", Port: 80, Open: true, Timestamp: time.Now()},
		{Host: "localhost", Port: 443, Open: true, Timestamp: time.Now()},
	}
}

func TestNew_CreatesDir(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "hist")
	_, err := history.New(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Fatal("expected directory to be created")
	}
}

func TestRecord_WritesFile(t *testing.T) {
	dir := t.TempDir()
	h, err := history.New(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	results := makeResults()
	if err := h.Record("localhost", results); err != nil {
		t.Fatalf("Record failed: %v", err)
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 file, got %d", len(entries))
	}
}

func TestRecord_ValidJSON(t *testing.T) {
	dir := t.TempDir()
	h, _ := history.New(dir)
	_ = h.Record("localhost", makeResults())

	entries, _ := os.ReadDir(dir)
	path := filepath.Join(dir, entries[0].Name())
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	var entry history.Entry
	if err := json.Unmarshal(data, &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if entry.Host != "localhost" {
		t.Errorf("expected host localhost, got %s", entry.Host)
	}
	if len(entry.Results) != 2 {
		t.Errorf("expected 2 results, got %d", len(entry.Results))
	}
}

func TestRecord_SanitizesHost(t *testing.T) {
	dir := t.TempDir()
	h, _ := history.New(dir)
	_ = h.Record("192.168.1.1:8080", makeResults())

	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Fatalf("expected 1 file, got %d", len(entries))
	}
	name := entries[0].Name()
	for _, c := range []string{":", "/"} {
		for _, ch := range name {
			if string(ch) == c {
				t.Errorf("filename contains unsafe char %q: %s", c, name)
			}
		}
	}
}
