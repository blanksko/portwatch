package snapshot

import "github.com/user/portwatch/internal/scanner"

// ChangeType describes what kind of change occurred on a port.
type ChangeType string

const (
	ChangeOpened ChangeType = "opened"
	ChangeClosed ChangeType = "closed"
)

// Change represents a detected difference between two snapshots.
type Change struct {
	Port   int        `json:"port"`
	Change ChangeType `json:"change"`
}

// Diff compares two snapshots and returns a list of changes.
// prev may be nil (first run), in which case all open ports are reported as opened.
func Diff(prev, curr *Snapshot) []Change {
	var changes []Change

	prevOpen := toPortSet(prev)
	currOpen := toPortSet(curr)

	for port := range currOpen {
		if !prevOpen[port] {
			changes = append(changes, Change{Port: port, Change: ChangeOpened})
		}
	}

	for port := range prevOpen {
		if !currOpen[port] {
			changes = append(changes, Change{Port: port, Change: ChangeClosed})
		}
	}

	return changes
}

func toPortSet(snap *Snapshot) map[int]struct{} {
	set := make(map[int]struct{})
	if snap == nil {
		return set
	}
	for _, r := range snap.Results {
		if r.Open {
			set[r.Port] = struct{}{}
		}
	}
	return set
}

// OpenPorts returns only the open ScanResults from a snapshot.
func OpenPorts(snap *Snapshot) []scanner.ScanResult {
	var open []scanner.ScanResult
	for _, r := range snap.Results {
		if r.Open {
			open = append(open, r)
		}
	}
	return open
}
