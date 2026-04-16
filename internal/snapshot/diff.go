package snapshot

import "github.com/user/portwatch/internal/scanner"

// DiffResult holds the sets of ports that were opened or closed
// between two snapshots.
type DiffResult struct {
	Opened []int
	Closed []int
}

// Diff compares prev and next scan results and returns the ports
// that were opened (present in next but not prev) and closed
// (present in prev but not next).
func Diff(prev, next []scanner.Result) DiffResult {
	prevSet := toPortSet(prev)
	nextSet := toPortSet(next)

	var opened, closed []int

	for p := range nextSet {
		if !prevSet[p] {
			opened = append(opened, p)
		}
	}
	for p := range prevSet {
		if !nextSet[p] {
			closed = append(closed, p)
		}
	}

	return DiffResult{Opened: opened, Closed: closed}
}

func toPortSet(results []scanner.Result) map[int]bool {
	s := make(map[int]bool, len(results))
	for _, r := range results {
		if r.Open {
			s[r.Port] = true
		}
	}
	return s
}

// OpenPorts returns only the open ports from a slice of Results.
func OpenPorts(results []scanner.Result) []int {
	var ports []int
	for _, r := range results {
		if r.Open {
			ports = append(ports, r.Port)
		}
	}
	return ports
}
