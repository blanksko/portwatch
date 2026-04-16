package scanner

import (
	"net"
	"strconv"
	"testing"
	"time"
)

// startTestServer starts a TCP listener on a random port and returns the port and a stop function.
func startTestServer(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test server: %v", err)
	}
	port, _ := strconv.Atoi(ln.Addr().(*net.TCPAddr).Port.String())
	_ = port
	actualPort := ln.Addr().(*net.TCPAddr).Port
	return actualPort, func() { ln.Close() }
}

func TestScan_OpenPort(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not start listener: %v", err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	s := New(500 * time.Millisecond)
	result, err := s.Scan("127.0.0.1", []int{port})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Ports) != 1 {
		t.Fatalf("expected 1 port result, got %d", len(result.Ports))
	}
	if !result.Ports[0].Open {
		t.Errorf("expected port %d to be open", port)
	}
}

func TestScan_ClosedPort(t *testing.T) {
	s := New(300 * time.Millisecond)
	// Port 1 is almost certainly closed in test environments.
	result, err := s.Scan("127.0.0.1", []int{1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Ports[0].Open {
		t.Errorf("expected port 1 to be closed")
	}
}

func TestScan_EmptyHost(t *testing.T) {
	s := New(300 * time.Millisecond)
	_, err := s.Scan("", []int{80})
	if err == nil {
		t.Error("expected error for empty host, got nil")
	}
}

func TestScanResult_Timestamp(t *testing.T) {
	s := New(300 * time.Millisecond)
	before := time.Now()
	result, _ := s.Scan("127.0.0.1", []int{})
	after := time.Now()
	if result.Timestamp.Before(before) || result.Timestamp.After(after) {
		t.Errorf("timestamp %v not in expected range", result.Timestamp)
	}
}
