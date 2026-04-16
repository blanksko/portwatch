package main

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.toml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempConfig: %v", err)
	}
	return p
}

func TestRun_NoConfig(t *testing.T) {
	// run() with a missing config should fall back to defaults and fail
	// validation because default hosts list is empty when no hosts configured.
	// We point to a non-existent path to trigger the fallback.
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"portwatch", "/nonexistent/portwatch.toml"}

	err := run()
	if err == nil {
		t.Fatal("expected error for empty hosts, got nil")
	}
}

func TestRun_InvalidConfig(t *testing.T) {
	p := writeTempConfig(t, `
[scan]
timeout = "2s"
`)
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"portwatch", p}

	err := run()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestRun_ValidConfig(t *testing.T) {
	dir := t.TempDir()
	snapshotPath := filepath.Join(dir, "snap.json")

	p := writeTempConfig(t, `
[scan]
timeout = "500ms"
ports = [19999]
snapshot_path = "`+snapshotPath+`"

[[hosts]]
address = "127.0.0.1"
`)
	origArgs := os.Args
	defer func() { os.Args = origArgs }()
	os.Args = []string{"portwatch", p}

	// First run: creates baseline snapshot.
	if err := run(); err != nil {
		t.Fatalf("first run: %v", err)
	}

	// Second run: diffs against baseline.
	if err := run(); err != nil {
		t.Fatalf("second run: %v", err)
	}
}
