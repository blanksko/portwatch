package config_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func writeTempConfig(t *testing.T, cfg *config.Config) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-cfg-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if err := json.NewEncoder(f).Encode(cfg); err != nil {
		t.Fatalf("encode config: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if len(cfg.Hosts) == 0 {
		t.Error("expected default hosts to be non-empty")
	}
	if len(cfg.Ports) == 0 {
		t.Error("expected default ports to be non-empty")
	}
	if cfg.TimeoutSeconds <= 0 {
		t.Error("expected positive default timeout")
	}
}

func TestLoad_Valid(t *testing.T) {
	cfg := &config.Config{
		Hosts:          []string{"192.168.1.1"},
		Ports:          []int{22, 443},
		SnapshotPath:   "/tmp/snap.json",
		TimeoutSeconds: 3,
	}
	path := writeTempConfig(t, cfg)
	loaded, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.Hosts[0] != "192.168.1.1" {
		t.Errorf("got host %q, want 192.168.1.1", loaded.Hosts[0])
	}
}

func TestLoad_NonExistent(t *testing.T) {
	_, err := config.Load("/no/such/file.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestValidate_NoHosts(t *testing.T) {
	cfg := &config.Config{Ports: []int{80}}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error when hosts is empty")
	}
}

func TestValidate_InvalidPort(t *testing.T) {
	cfg := &config.Config{
		Hosts: []string{"localhost"},
		Ports: []int{0},
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for port 0")
	}
}

func TestValidate_DefaultsApplied(t *testing.T) {
	cfg := &config.Config{
		Hosts: []string{"localhost"},
		Ports: []int{80},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.TimeoutSeconds != 2 {
		t.Errorf("expected default timeout 2, got %d", cfg.TimeoutSeconds)
	}
	if cfg.SnapshotPath == "" {
		t.Error("expected non-empty default snapshot path")
	}
}
