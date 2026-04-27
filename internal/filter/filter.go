// Package filter provides port filtering utilities for portwatch.
package filter

import "github.com/user/portwatch/internal/scanner"

// Rule defines criteria for including or excluding ports from results.
type Rule struct {
	Ports    []int
	Protocol string
}

// Filter applies inclusion and exclusion rules to scan results.
type Filter struct {
	include map[int]struct{}
	exclude map[int]struct{}
}

// New creates a Filter. include and exclude are port lists; empty include means allow all.
func New(include, exclude []int) *Filter {
	f := &Filter{
		include: toSet(include),
		exclude: toSet(exclude),
	}
	return f
}

// Apply returns only the scan results that pass the filter rules.
func (f *Filter) Apply(results []scanner.Result) []scanner.Result {
	out := make([]scanner.Result, 0, len(results))
	for _, r := range results {
		if f.allow(r.Port) {
			out = append(out, r)
		}
	}
	return out
}

// Count returns the number of results that would pass the filter rules
// without allocating a new slice.
func (f *Filter) Count(results []scanner.Result) int {
	n := 0
	for _, r := range results {
		if f.allow(r.Port) {
			n++
		}
	}
	return n
}

func (f *Filter) allow(port int) bool {
	if _, excluded := f.exclude[port]; excluded {
		return false
	}
	if len(f.include) == 0 {
		return true
	}
	_, included := f.include[port]
	return included
}

func toSet(ports []int) map[int]struct{} {
	s := make(map[int]struct{}, len(ports))
	for _, p := range ports {
		s[p] = struct{}{}
	}
	return s
}
