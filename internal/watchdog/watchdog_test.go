package watchdog_test

import (
	"bytes"
	"log"
	"testing"

	"github.com/example/portwatch/internal/scanner"
	"github.com	/example/portwatch/internal/watchdog"
)

func silentLogger() *log.Logger {
	return log.New(&bytes.Buffer{}, "", 0)
}

func TestNew_MissingScanner(t *testing.T) {
	_, err := watchdog.New(watchdog.Config{
		Logger: silentLogger(),
	})
	if err == nil {
		t.Fatal("expected error for nil scanner")
	}
}

func TestNew_MissingLogger(t *testing.T) {
	s, _ := scanner.New(scanner.Config{Timeout: 0})
	_, err := watchdog.New(watchdog.Config{
		Scanner: s,
	})
	if err == nil {
		t.Fatal("expected error for nil logger")
	}
}

func TestRun_NoHostsSucceeds(t *testing.T) {
	s, err := scanner.New(scanner.Config{Timeout: 0})
	if err != nil {
		t.Fatalf("scanner.New: %v", err)
	}
	w, err := watchdog.New(watchdog.Config{
		Scanner: s,
		Logger:  silentLogger(),
	})
	if err != nil {
		t.Fatalf("watchdog.New: %v", err)
	}
	if err := w.Run(nil, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRun_ReturnsNoErrorOnEmptyPorts(t *testing.T) {
	s, _ := scanner.New(scanner.Config{Timeout: 0})
	w, _ := watchdog.New(watchdog.Config{
		Scanner: s,
		Logger:  silentLogger(),
	})
	if err := w.Run([]string{"127.0.0.1"}, []int{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
