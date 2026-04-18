package baseline

import (
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

// Deviation describes ports that deviate from the approved baseline.
type Deviation struct {
	Host     string
	Extra    []int // open but not in baseline
	Missing  []int // in baseline but no longer open
}

// Compare checks live results against a stored baseline.
// If no baseline exists for the host, Compare returns nil.
func Compare(b *Baseline, live []scanner.Result) *Deviation {
	if b == nil {
		return nil
	}

	d := snapshot.Diff(b.Ports, live)
	if len(d.Opened) == 0 && len(d.Closed) == 0 {
		return nil
	}

	dev := &Deviation{Host: b.Host}
	for _, p := range d.Opened {
		dev.Extra = append(dev.Extra, p)
	}
	for _, p := range d.Closed {
		dev.Missing = append(dev.Missing, p)
	}
	return dev
}
