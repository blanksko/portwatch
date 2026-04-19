package grouping_test

import (
	"testing"

	"github.com/example/portwatch/internal/grouping"
)

func TestAdd_And_Get(t *testing.T) {
	r := grouping.New()
	r.Add("web", "10.0.0.1", "10.0.0.2")
	hosts, ok := r.Get("web")
	if !ok {
		t.Fatal("expected group 'web' to exist")
	}
	if len(hosts) != 2 {
		t.Fatalf("expected 2 hosts, got %d", len(hosts))
	}
}

func TestGet_UnknownGroup(t *testing.T) {
	r := grouping.New()
	_, ok := r.Get("unknown")
	if ok {
		t.Fatal("expected false for unknown group")
	}
}

func TestAdd_MergesHosts(t *testing.T) {
	r := grouping.New()
	r.Add("db", "10.0.1.1")
	r.Add("db", "10.0.1.2")
	hosts, _ := r.Get("db")
	if len(hosts) != 2 {
		t.Fatalf("expected 2 hosts after merge, got %d", len(hosts))
	}
}

func TestDelete_RemovesGroup(t *testing.T) {
	r := grouping.New()
	r.Add("tmp", "192.168.1.1")
	r.Delete("tmp")
	_, ok := r.Get("tmp")
	if ok {
		t.Fatal("expected group to be deleted")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	r := grouping.New()
	r.Add("a", "h1")
	r.Add("b", "h2", "h3")
	all := r.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(all))
	}
}

func TestHosts_Unique(t *testing.T) {
	r := grouping.New()
	r.Add("x", "shared", "only-x")
	r.Add("y", "shared", "only-y")
	hosts := r.Hosts()
	if len(hosts) != 3 {
		t.Fatalf("expected 3 unique hosts, got %d", len(hosts))
	}
}
