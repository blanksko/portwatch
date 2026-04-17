// Package watchdog ties together scanning, diffing, alerting, and history
// recording into a single reusable job that can be driven by the scheduler.
package watchdog

import (
	"fmt"
	"log"

	"github.com/example/portwatch/internal/history"
	"github.com/example/portwatch/internal/metrics"
	"github.com/example/portwatch/internal/notifier"
	"github.com/example/portwatch/internal/scanner"
	"github.com/example/portwatch/internal/snapshot"
)

// Watchdog performs a full scan cycle for a set of hosts and ports.
type Watchdog struct {
	scanner  *scanner.Scanner
	snap     *snapshot.Snapshot
	notifier *notifier.Notifier
	history  *history.History
	metrics  *metrics.Metrics
	logger   *log.Logger
}

// Config holds the dependencies needed to build a Watchdog.
type Config struct {
	Scanner  *scanner.Scanner
	Snapshot *snapshot.Snapshot
	Notifier *notifier.Notifier
	History  *history.History
	Metrics  *metrics.Metrics
	Logger   *log.Logger
}

// New creates a Watchdog from the provided Config.
func New(cfg Config) (*Watchdog, error) {
	if cfg.Scanner == nil {
		return nil, fmt.Errorf("watchdog: scanner is required")
	}
	if cfg.Logger == nil {
		return nil, fmt.Errorf("watchdog: logger is required")
	}
	return &Watchdog{
		scanner:  cfg.Scanner,
		snap:     cfg.Snapshot,
		notifier: cfg.Notifier,
		history:  cfg.History,
		metrics:  cfg.Metrics,
		logger:   cfg.Logger,
	}, nil
}

// Run executes one full scan cycle: scan → diff → notify → record.
func (w *Watchdog) Run(hosts []string, ports []int) error {
	results, err := w.scanner.Scan(hosts, ports)
	if err != nil {
		return fmt.Errorf("watchdog: scan failed: %w", err)
	}

	if w.metrics != nil {
		w.metrics.Record(results)
	}

	if w.snap != nil {
		prev, _ := w.snap.Load()
		diff := snapshot.Diff(prev, results)
		if w.notifier != nil {
			if err := w.notifier.Notify(diff); err != nil {
				w.logger.Printf("watchdog: notify error: %v", err)
			}
		}
		if err := w.snap.Save(results); err != nil {
			w.logger.Printf("watchdog: snapshot save error: %v", err)
		}
	}

	if w.history != nil {
		for _, r := range results {
			if err := w.history.Record(r.Host, results); err != nil {
				w.logger.Printf("watchdog: history record error: %v", err)
			}
			break // record once per cycle keyed to first host; history handles per-host internally
		}
	}

	w.logger.Printf("watchdog: cycle complete, %d results", len(results))
	return nil
}
