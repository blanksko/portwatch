package shadow

import (
	"github.com/user/portwatch/internal/scanner"
)

// Delta describes ports that differ between a shadow record and a live scan.
type Delta struct {
	Host    string
	Added   []int // present in live, absent in shadow
	Removed []int // present in shadow, absent in live
}

// HasChanges reports whether the delta contains any differences.
func (d Delta) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0
}

// Compare returns the delta between a previously recorded shadow snapshot
// and the current live results for the same host.
func Compare(rec Record, live []scanner.Result) Delta {
	shadowPorts := toSet(rec.Results)
	livePorts := toSet(live)

	d := Delta{Host: rec.Host}

	for p := range livePorts {
		if !shadowPorts[p] {
			d.Added = append(d.Added, p)
		}
	}
	for p := range shadowPorts {
		if !livePorts[p] {
			d.Removed = append(d.Removed, p)
		}
	}
	return d
}

func toSet(results []scanner.Result) map[int]bool {
	s := make(map[int]bool, len(results))
	for _, r := range results {
		if r.Open {
			s[r.Port] = true
		}
	}
	return s
}
