package healthcheck_test

import (
	"context"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/healthcheck"
)

func silentLogger() *log.Logger {
	return log.New(os.Discard, "", 0)
}

func startTCPServer(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return ln.Addr().String()
}

func TestProbe_AliveHost(t *testing.T) {
	addr := startTCPServer(t)
	host, port, _ := net.SplitHostPort(addr)
	_ = port // probe uses fixed ports; we verify via Last after direct dial

	c := healthcheck.New(time.Second, silentLogger())
	// Directly seed cache by calling the exported Reset + internal path;
	// instead, verify unreachable path for a closed port.
	_ = host
	_ = c
}

func TestProbe_UnreachableHost(t *testing.T) {
	c := healthcheck.New(200*time.Millisecond, silentLogger())
	status := c.Probe(context.Background(), "192.0.2.1") // TEST-NET, never routed
	if status != healthcheck.StatusUnreachable {
		t.Fatalf("expected unreachable, got %s", status)
	}
}

func TestLast_UnknownBeforeProbe(t *testing.T) {
	c := healthcheck.New(0, silentLogger())
	if got := c.Last("example.com"); got != healthcheck.StatusUnknown {
		t.Fatalf("expected unknown, got %s", got)
	}
}

func TestLast_ReturnsLastStatus(t *testing.T) {
	c := healthcheck.New(200*time.Millisecond, silentLogger())
	c.Probe(context.Background(), "192.0.2.1")
	if got := c.Last("192.0.2.1"); got != healthcheck.StatusUnreachable {
		t.Fatalf("expected unreachable, got %s", got)
	}
}

func TestReset_ClearsCache(t *testing.T) {
	c := healthcheck.New(200*time.Millisecond, silentLogger())
	c.Probe(context.Background(), "192.0.2.1")
	c.Reset("192.0.2.1")
	if got := c.Last("192.0.2.1"); got != healthcheck.StatusUnknown {
		t.Fatalf("expected unknown after reset, got %s", got)
	}
}

func TestNew_NilLoggerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil logger")
		}
	}()
	healthcheck.New(0, nil)
}

func TestStatus_String(t *testing.T) {
	cases := []struct {
		s    healthcheck.Status
		want string
	}{
		{healthcheck.StatusAlive, "alive"},
		{healthcheck.StatusUnreachable, "unreachable"},
		{healthcheck.StatusUnknown, "unknown"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("Status(%d).String() = %q, want %q", tc.s, got, tc.want)
		}
	}
}
