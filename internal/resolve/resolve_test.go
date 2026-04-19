package resolve_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/portwatch/internal/resolve"
)

func TestResolve_IPPassthrough(t *testing.T) {
	r := resolve.New(2 * time.Second)
	res, err := r.Resolve("127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Addresses) != 1 || res.Addresses[0] != "127.0.0.1" {
		t.Errorf("expected [127.0.0.1], got %v", res.Addresses)
	}
}

func TestResolve_Localhost(t *testing.T) {
	r := resolve.New(2 * time.Second)
	res, err := r.Resolve("localhost")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Addresses) == 0 {
		t.Error("expected at least one address for localhost")
	}
	if res.ResolvedAt.IsZero() {
		t.Error("expected non-zero ResolvedAt")
	}
}

func TestPrimary_IP(t *testing.T) {
	r := resolve.New(2 * time.Second)
	ip, err := r.Primary("127.0.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ip != "127.0.0.1" {
		t.Errorf("expected 127.0.0.1, got %s", ip)
	}
}

func TestResolve_InvalidHost(t *testing.T) {
	r := resolve.New(2 * time.Second)
	_, err := r.Resolve("this.host.does.not.exist.invalid")
	if err == nil {
		t.Fatal("expected error for invalid host")
	}
	if !strings.Contains(err.Error(), "resolve") {
		t.Errorf("expected 'resolve' in error message, got: %v", err)
	}
}

func TestNew_DefaultTimeout(t *testing.T) {
	r := resolve.New(0)
	if r == nil {
		t.Fatal("expected non-nil resolver")
	}
}
