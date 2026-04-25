package pipeline

import (
	"context"
	"log"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/tagger"
)

// Builder assembles a Pipeline from high-level options.
type Builder struct {
	logger  *log.Logger
	filter  *filter.Filter
	tagger  *tagger.Tagger
	extra   []Stage
}

// NewBuilder returns a Builder backed by the given logger.
func NewBuilder(logger *log.Logger) *Builder {
	if logger == nil {
		panic("pipeline: builder logger must not be nil")
	}
	return &Builder{logger: logger}
}

// WithFilter attaches a port filter stage.
func (b *Builder) WithFilter(f *filter.Filter) *Builder {
	b.filter = f
	return b
}

// WithTagger attaches a port-tagging stage.
func (b *Builder) WithTagger(t *tagger.Tagger) *Builder {
	b.tagger = t
	return b
}

// WithStage appends a custom stage.
func (b *Builder) WithStage(s Stage) *Builder {
	b.extra = append(b.extra, s)
	return b
}

// Build constructs the Pipeline in a sensible default order:
// filter → tagger → custom stages.
func (b *Builder) Build() *Pipeline {
	p := New(b.logger)

	if b.filter != nil {
		f := b.filter
		p.Use(Stage{
			Name: "filter",
			Handler: func(_ context.Context, results []scanner.Result, next func([]scanner.Result) error) error {
				return next(f.Apply(results))
			},
		})
	}

	if b.tagger != nil {
		t := b.tagger
		p.Use(Stage{
			Name: "tagger",
			Handler: func(_ context.Context, results []scanner.Result, next func([]scanner.Result) error) error {
				for i := range results {
					results[i] = t.Tag(results[i])
				}
				return next(results)
			},
		})
	}

	for _, s := range b.extra {
		p.Use(s)
	}

	return p
}
