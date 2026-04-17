// Package metrics tracks scan run statistics such as open port counts
// and scan durations for reporting and diagnostics.
package metrics

import (
	"sync"
	"time"
)

// Snapshot holds a point-in-time view of collected metrics.
type Snapshot struct {
	TotalScans   int
	OpenPorts    int
	LastScanAt   time.Time
	LastDuration time.Duration
}

// Collector accumulates scan metrics in memory.
type Collector struct {
	mu           sync.Mutex
	totalScans   int
	openPorts    int
	lastScanAt   time.Time
	lastDuration time.Duration
}

// New returns a new Collector.
func New() *Collector {
	return &Collector{}
}

// Record records the result of a single scan run.
func (c *Collector) Record(openPorts int, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.totalScans++
	c.openPorts = openPorts
	c.lastScanAt = time.Now()
	c.lastDuration = duration
}

// Snapshot returns a copy of the current metrics.
func (c *Collector) Snapshot() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Snapshot{
		TotalScans:   c.totalScans,
		OpenPorts:    c.openPorts,
		LastScanAt:   c.lastScanAt,
		LastDuration: c.lastDuration,
	}
}

// Reset clears all accumulated metrics.
func (c *Collector) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	*c = Collector{}
}
