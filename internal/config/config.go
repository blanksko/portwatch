// Package config handles loading and validating portwatch configuration.
package config

import (
	"encoding/json"
	"errors"
	"os"
)

// Config holds the top-level portwatch configuration.
type Config struct {
	// Hosts is the list of hosts to scan.
	Hosts []string `json:"hosts"`
	// Ports is the list of ports to check on each host.
	Ports []int `json:"ports"`
	// SnapshotPath is the file path used to persist snapshots.
	SnapshotPath string `json:"snapshot_path"`
	// TimeoutSeconds is the per-port dial timeout.
	TimeoutSeconds int `json:"timeout_seconds"`
}

// Default returns a Config populated with sensible defaults.
func Default() *Config {
	return &Config{
		Hosts:          []string{"localhost"},
		Ports:          []int{22, 80, 443, 8080},
		SnapshotPath:   "portwatch_snapshot.json",
		TimeoutSeconds: 2,
	}
}

// Load reads a JSON config file from path and returns a validated Config.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Validate checks that required fields are present and values are in range.
func (c *Config) Validate() error {
	if len(c.Hosts) == 0 {
		return errors.New("config: at least one host is required")
	}
	if len(c.Ports) == 0 {
		return errors.New("config: at least one port is required")
	}
	for _, p := range c.Ports {
		if p < 1 || p > 65535 {
			return errors.New("config: port out of range (1-65535)")
		}
	}
	if c.TimeoutSeconds <= 0 {
		c.TimeoutSeconds = 2
	}
	if c.SnapshotPath == "" {
		c.SnapshotPath = "portwatch_snapshot.json"
	}
	return nil
}
