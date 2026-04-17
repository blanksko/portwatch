package watchdog_test

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/config"
	"github.com/example/portwatch/internal/watchdog"
)

func validConfig(t *testing.T) *config.Config {
	t.Helper()
	cfg := config.Default()
	cfg.Hosts = []string{"127.0.0.1"}
	cfg.Ports = []int{80}
	cfg.Timeout = 100 * time.Millisecond
	cfg.SnapshotPath = t.TempDir() + "/snap.json"
	cfg.HistoryDir = t.TempDir()
	return cfg
}

func TestFromConfig_ValidConfig(t *testing.T) {
	w, err := watchdog.FromConfig(validConfig(t), silentLogger())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if w == nil {
		t.Fatal("expected non-nil watchdog")
	}
}

func TestFromConfig_NilConfig(t *testing.T) {
	_, err := watchdog.FromConfig(nil, silentLogger())
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestFromConfig_RunCycle(t *testing.T) {
	w, err := watchdog.FromConfig(validConfig(t), silentLogger())
	if err != nil {
		t.Fatalf("FromConfig: %v", err)
	}
	if err := w.Run([]string{"127.0.0.1"}, []int{1, 2, 3}); err != nil {
		t.Fatalf("Run: %v", err)
	}
}
