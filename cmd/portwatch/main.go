package main

import (
	"fmt"
	"os"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfgPath := "portwatch.toml"
	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warn: could not load config (%v), using defaults\n", err)
		cfg = config.Default()
	}

	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	sc := scanner.New(cfg.Timeout)

	var allResults []scanner.Result
	for _, host := range cfg.Hosts {
		results, err := sc.Scan(host, cfg.Ports)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: scan failed for %s: %v\n", host, err)
			continue
		}
		allResults = append(allResults, results...)
	}

	prev, err := snapshot.Load(cfg.SnapshotPath)
	if err != nil {
		fmt.Println("no previous snapshot found, saving baseline")
		return snapshot.Save(cfg.SnapshotPath, allResults)
	}

	diff := snapshot.Diff(prev, allResults)

	al := alert.New(os.Stdout)
	if err := al.Notify(diff); err != nil {
		return fmt.Errorf("alert notify: %w", err)
	}

	return snapshot.Save(cfg.SnapshotPath, allResults)
}
