package tagger_test

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tagger"
)

func makeResult(host string, port int) scanner.Result {
	return scanner.Result{Host: host, Port: port, Open: true}
}

func TestTag_WellKnown(t *testing.T) {
	tg := tagger.Default()
	tags := tg.Tag(80)
	if !contains(tags, "well-known") {
		t.Errorf("expected well-known tag for port 80, got %v", tags)
	}
	if !contains(tags, "http") {
		t.Errorf("expected http tag for port 80, got %v", tags)
	}
}

func TestTag_SSH(t *testing.T) {
	tg := tagger.Default()
	tags := tg.Tag(22)
	if !contains(tags, "ssh") {
		t.Errorf("expected ssh tag for port 22, got %v", tags)
	}
}

func TestTag_Dynamic(t *testing.T) {
	tg := tagger.Default()
	tags := tg.Tag(60000)
	if !contains(tags, "dynamic") {
		t.Errorf("expected dynamic tag for port 60000, got %v", tags)
	}
}

func TestTag_NoMatch(t *testing.T) {
	tg := tagger.New([]tagger.Rule{})
	tags := tg.Tag(80)
	if len(tags) != 0 {
		t.Errorf("expected no tags, got %v", tags)
	}
}

func TestAnnotate_AssignsTags(t *testing.T) {
	tg := tagger.Default()
	results := []scanner.Result{
		makeResult("localhost", 443),
		makeResult("localhost", 50000),
	}
	annotated := tg.Annotate(results)
	if len(annotated) != 2 {
		t.Fatalf("expected 2 annotated results, got %d", len(annotated))
	}
	if !contains(annotated[0].Tags, "https") {
		t.Errorf("expected https tag on port 443, got %v", annotated[0].Tags)
	}
	if !contains(annotated[1].Tags, "dynamic") {
		t.Errorf("expected dynamic tag on port 50000, got %v", annotated[1].Tags)
	}
}

func TestAnnotatedResult_String(t *testing.T) {
	tg := tagger.Default()
	annotated := tg.Annotate([]scanner.Result{makeResult("127.0.0.1", 22)})
	s := annotated[0].String()
	if s == "" {
		t.Error("expected non-empty string representation")
	}
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
