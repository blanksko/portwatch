// Package scheduler provides periodic scan scheduling for portwatch.
package scheduler

import (
	"context"
	"log"
	"time"
)

// Job is a function executed on each tick.
type Job func(ctx context.Context) error

// Scheduler runs a Job at a fixed interval.
type Scheduler struct {
	interval time.Duration
	job      Job
	logger   *log.Logger
}

// New creates a Scheduler with the given interval and job.
func New(interval time.Duration, job Job, logger *log.Logger) *Scheduler {
	if logger == nil {
		logger = log.Default()
	}
	return &Scheduler{
		interval: interval,
		job:      job,
		logger:   logger,
	}
}

// Run starts the scheduler loop, blocking until ctx is cancelled.
// The job is executed immediately on start, then on every tick.
func (s *Scheduler) Run(ctx context.Context) {
	s.logger.Printf("scheduler: starting with interval %s", s.interval)

	if err := s.job(ctx); err != nil {
		s.logger.Printf("scheduler: job error: %v", err)
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.logger.Println("scheduler: stopped")
			return
		case <-ticker.C:
			if err := s.job(ctx); err != nil {
				s.logger.Printf("scheduler: job error: %v", err)
			}
		}
	}
}
