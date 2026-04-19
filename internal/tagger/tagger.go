// Package tagger assigns automatic tags to scan results based on port ranges and protocols.
package tagger

import (
	"fmt"

	"github.com/user/portwatch/internal/scanner"
)

// Rule defines a tagging rule for a port range.
type Rule struct {
	Min  int
	Max  int
	Tag  string
}

// Tagger applies tag rules to scan results.
type Tagger struct {
	rules []Rule
}

// Default returns a Tagger with common well-known port range rules.
func Default() *Tagger {
	return New([]Rule{
		{Min: 1, Max: 1023, Tag: "well-known"},
		{Min: 1024, Max: 49151, Tag: "registered"},
		{Min: 49152, Max: 65535, Tag: "dynamic"},
		{Min: 80, Max: 80, Tag: "http"},
		{Min: 443, Max: 443, Tag: "https"},
		{Min: 22, Max: 22, Tag: "ssh"},
		{Min: 3306, Max: 3306, Tag: "mysql"},
		{Min: 5432, Max: 5432, Tag: "postgres"},
	})
}

// New creates a Tagger with the provided rules.
func New(rules []Rule) *Tagger {
	return &Tagger{rules: rules}
}

// Tag returns all matching tags for a given port number.
func (t *Tagger) Tag(port int) []string {
	seen := map[string]struct{}{}
	var tags []string
	for _, r := range t.rules {
		if port >= r.Min && port <= r.Max {
			if _, ok := seen[r.Tag]; !ok {
				tags = append(tags, r.Tag)
				seen[r.Tag] = struct{}{}
			}
		}
	}
	return tags
}

// Annotate enriches each ScanResult's label with matched tags.
func (t *Tagger) Annotate(results []scanner.Result) []AnnotatedResult {
	out := make([]AnnotatedResult, 0, len(results))
	for _, r := range results {
		out = append(out, AnnotatedResult{
			Result: r,
			Tags:   t.Tag(r.Port),
		})
	}
	return out
}

// AnnotatedResult wraps a scanner.Result with computed tags.
type AnnotatedResult struct {
	scanner.Result
	Tags []string
}

// String returns a human-readable representation.
func (a AnnotatedResult) String() string {
	return fmt.Sprintf("%s:%d %v", a.Host, a.Port, a.Tags)
}
