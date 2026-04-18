package metrics

import (
	"encoding/json"
	"io"
	"time"
)

// Snapshot holds an exported view of current metrics.
type Snapshot struct {
	TotalScans   int64     `json:"total_scans"`
	OpenPorts    int       `json:"open_ports"`
	ClosedPorts  int       `json:"closed_ports"`
	LastScanTime time.Time `json:"last_scan_time"`
	ExportedAt   time.Time `json:"exported_at"`
}

// Exporter writes metrics snapshots as JSON to a writer.
type Exporter struct {
	m *Metrics
}

// NewExporter creates an Exporter backed by the given Metrics.
func NewExporter(m *Metrics) *Exporter {
	if m == nil {
		panic("metrics: NewExporter requires non-nil Metrics")
	}
	return &Exporter{m: m}
}

// Export serialises the current metrics state as JSON into w.
func (e *Exporter) Export(w io.Writer) error {
	e.m.mu.Lock()
	snap := Snapshot{
		TotalScans:   e.m.totalScans,
		OpenPorts:    e.m.openPorts,
		ClosedPorts:  e.m.closedPorts,
		LastScanTime: e.m.lastScanTime,
		ExportedAt:   time.Now().UTC(),
	}
	e.m.mu.Unlock()

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(snap)
}
