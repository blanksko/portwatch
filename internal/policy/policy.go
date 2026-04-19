// Package policy evaluates whether a port change should trigger an alert
// based on configurable rules such as port ranges and protocols.
package policy

import "github.com/user/portwatch/internal/snapshot"

// Rule defines a single policy rule.
type Rule struct {
	Ports    []int  // empty means all ports
	Action   string // "allow" or "deny"
}

// Policy holds an ordered list of rules evaluated top-down.
type Policy struct {
	rules []Rule
}

// New returns a Policy with the given rules.
func New(rules []Rule) *Policy {
	return &Policy{rules: rules}
}

// Default returns a Policy that allows all port changes.
func Default() *Policy {
	return &Policy{}
}

// Allow reports whether the opened/closed port change should produce an alert.
// Returns true when no rules match (default allow) or the matching rule action is "allow".
func (p *Policy) Allow(port int) bool {
	for _, r := range p.rules {
		if matches(r, port) {
			return r.Action == "allow"
		}
	}
	return true
}

// Filter removes diff entries that are suppressed by policy.
func (p *Policy) Filter(d snapshot.Diff) snapshot.Diff {
	if len(p.rules) == 0 {
		return d
	}
	out := snapshot.Diff{}
	for _, r := range d.Opened {
		if p.Allow(r.Port) {
			out.Opened = append(out.Opened, r)
		}
	}
	for _, r := range d.Closed {
		if p.Allow(r.Port) {
			out.Closed = append(out.Closed, r)
		}
	}
	return out
}

func matches(r Rule, port int) bool {
	if len(r.Ports) == 0 {
		return true
	}
	for _, p := range r.Ports {
		if p == port {
			return true
		}
	}
	return false
}
