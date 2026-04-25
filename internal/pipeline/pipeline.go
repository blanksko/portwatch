// Package pipeline chains scan middleware stages into an ordered execution pipeline.
package pipeline

import (
	"context"
	"fmt"
	"log"

	"github.com/user/portwatch/internal/scanner"
)

// Stage is a single step in the scan pipeline.
type Stage struct {
	Name    string
	Handler func(ctx context.Context, results []scanner.Result, next func([]scanner.Result) error) error
}

// Pipeline executes a sequence of stages around a terminal scan job.
type Pipeline struct {
	stages []Stage
	logger *log.Logger
}

// New returns a Pipeline with the provided logger.
func New(logger *log.Logger) *Pipeline {
	if logger == nil {
		panic("pipeline: logger must not be nil")
	}
	return &Pipeline{logger: logger}
}

// Use appends a stage to the pipeline.
func (p *Pipeline) Use(s Stage) {
	p.stages = append(p.stages, s)
}

// Run executes the pipeline, calling each stage in order and finally invoking
// terminal with the (possibly modified) results.
func (p *Pipeline) Run(ctx context.Context, results []scanner.Result, terminal func([]scanner.Result) error) error {
	return p.run(ctx, results, 0, terminal)
}

func (p *Pipeline) run(ctx context.Context, results []scanner.Result, idx int, terminal func([]scanner.Result) error) error {
	if idx >= len(p.stages) {
		return terminal(results)
	}
	s := p.stages[idx]
	p.logger.Printf("pipeline: executing stage %q", s.Name)
	err := s.Handler(ctx, results, func(next []scanner.Result) error {
		return p.run(ctx, next, idx+1, terminal)
	})
	if err != nil {
		return fmt.Errorf("pipeline stage %q: %w", s.Name, err)
	}
	return nil
}
