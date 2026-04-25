package pipeline_test

import (
	"context"
	"testing"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tagger"
)

func TestBuilder_NoStages_PassesThrough(t *testing.T) {
	p := pipeline.NewBuilder(silentLogger()).Build()
	results := makeResults(22, 80)
	var got []scanner.Result
	if err := p.Run(context.Background(), results, func(r []scanner.Result) error {
		got = r
		return nil
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 results, got %d", len(got))
	}
}

func TestBuilder_WithFilter_RemovesPorts(t *testing.T) {
	f := filter.New(filter.Options{Exclude: []int{22}})
	p := pipeline.NewBuilder(silentLogger()).WithFilter(f).Build()
	results := makeResults(22, 80)
	var got []scanner.Result
	_ = p.Run(context.Background(), results, func(r []scanner.Result) error {
		got = r
		return nil
	})
	for _, r := range got {
		if r.Port == 22 {
			t.Fatal("port 22 should have been filtered out")
		}
	}
}

func TestBuilder_WithTagger_TagsResults(t *testing.T) {
	tgr := tagger.Default()
	p := pipeline.NewBuilder(silentLogger()).WithTagger(tgr).Build()
	results := makeResults(22)
	var got []scanner.Result
	_ = p.Run(context.Background(), results, func(r []scanner.Result) error {
		got = r
		return nil
	})
	if len(got) == 0 {
		t.Fatal("expected at least one result")
	}
	if got[0].Tag == "" {
		t.Fatal("expected port 22 to be tagged")
	}
}

func TestBuilder_NilLoggerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil logger")
		}
	}()
	pipeline.NewBuilder(nil)
}
