package tags_test

import (
	"testing"

	"github.com/yourorg/portwatch/internal/tags"
)

func makeTags() *tags.Tags {
	return tags.New(map[string][]string{
		"192.168.1.1": {"production", "web"},
		"192.168.1.2": {"staging"},
		"10.0.0.1":    {"production", "db"},
	})
}

func TestGet_KnownHost(t *testing.T) {
	tg := makeTags()
	got := tg.Get("192.168.1.1")
	if len(got) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(got))
	}
}

func TestGet_UnknownHost(t *testing.T) {
	tg := makeTags()
	got := tg.Get("1.2.3.4")
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}

func TestHas_MatchingTag(t *testing.T) {
	tg := makeTags()
	if !tg.Has("192.168.1.1", "web") {
		t.Fatal("expected host to have tag 'web'")
	}
}

func TestHas_MissingTag(t *testing.T) {
	tg := makeTags()
	if tg.Has("192.168.1.2", "production") {
		t.Fatal("expected host NOT to have tag 'production'")
	}
}

func TestHosts_ReturnsTaggedHosts(t *testing.T) {
	tg := makeTags()
	hosts := tg.Hosts("production")
	if len(hosts) != 2 {
		t.Fatalf("expected 2 production hosts, got %d", len(hosts))
	}
}

func TestHosts_NoMatch(t *testing.T) {
	tg := makeTags()
	hosts := tg.Hosts("nonexistent")
	if len(hosts) != 0 {
		t.Fatalf("expected 0 hosts, got %d", len(hosts))
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	tg := makeTags()
	all := tg.All()
	if len(all) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(all))
	}
	// mutating the copy should not affect internal state
	delete(all, "10.0.0.1")
	if len(tg.All()) != 3 {
		t.Fatal("All() should return an independent copy")
	}
}

func TestNew_TrimsWhitespace(t *testing.T) {
	tg := tags.New(map[string][]string{
		" host1 ": {" prod ", "web"},
	})
	if !tg.Has("host1", "prod") {
		t.Fatal("expected whitespace to be trimmed from host and tag")
	}
}
