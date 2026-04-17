package watchdog

import (
	"fmt"
	"log"
	"os"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/config"
	"github.com/example/portwatch/internal/history"
	"github.com/example/portwatch/internal/metrics"
	"github.com/example/portwatch/internal/notifier"
	"github.com/example/portwatch/internal/reporter"
	"github.com/example/portwatch/internal/scanner"
	"github.com/example/portwatch/internal/snapshot"
)

// FromConfig builds a ready-to-use Watchdog from application config.
func FromConfig(cfg *config.Config, logger *log.Logger) (*Watchdog, error) {
	if cfg == nil {
		return nil, fmt.Errorf("watchdog: config is nil")
	}

	s, err := scanner.New(scanner.Config{Timeout: cfg.Timeout})
	if err != nil {
		return nil, fmt.Errorf("watchdog: scanner: %w", err)
	}

	snap := snapshot.New(cfg.SnapshotPath)

	alertN, err := alert.New(alert.Config{Writer: os.Stdout})
	if err != nil {
		return nil, fmt.Errorf("watchdog: alert: %w", err)
	}
	rep, err := reporter.New(reporter.Config{Writer: os.Stdout})
	if err != nil {
		return nil, fmt.Errorf("watchdog: reporter: %w", err)
	}
	multi := notifier.NewMulti(alertN, rep)
	n, err := notifier.New(notifier.Config{Delegate: multi})
	if err != nil {
		return nil, fmt.Errorf("watchdog: notifier: %w", err)
	}

	h, err := history.New(cfg.HistoryDir)
	if err != nil {
		return nil, fmt.Errorf("watchdog: history: %w", err)
	}

	m := metrics.New()

	return New(Config{
		Scanner:  s,
		Snapshot: snap,
		Notifier: n,
		History:  h,
		Metrics:  m,
		Logger:   logger,
	})
}
