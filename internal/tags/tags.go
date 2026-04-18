// Package tags provides host tagging support for portwatch,
// allowing users to group and label monitored hosts.
package tags

import "strings"

// Tags holds a mapping of host to its associated labels.
type Tags struct {
	hostTags map[string][]string
}

// New creates a new Tags instance from a map of host -> tag list.
func New(hostTags map[string][]string) *Tags {
	normalized := make(map[string][]string, len(hostTags))
	for host, tags := range hostTags {
		copy := make([]string, len(tags))
		for i, t := range tags {
			copy[i] = strings.TrimSpace(t)
		}
		normalized[strings.TrimSpace(host)] = copy
	}
	return &Tags{hostTags: normalized}
}

// Get returns the tags associated with a host.
// Returns an empty slice if the host has no tags.
func (t *Tags) Get(host string) []string {
	tags, ok := t.hostTags[strings.TrimSpace(host)]
	if !ok {
		return []string{}
	}
	return tags
}

// Has reports whether the given host has the specified tag.
func (t *Tags) Has(host, tag string) bool {
	for _, tg := range t.Get(host) {
		if tg == strings.TrimSpace(tag) {
			return true
		}
	}
	return false
}

// Hosts returns all hosts that carry the given tag.
func (t *Tags) Hosts(tag string) []string {
	tag = strings.TrimSpace(tag)
	var result []string
	for host, tags := range t.hostTags {
		for _, tg := range tags {
			if tg == tag {
				result = append(result, host)
				break
			}
		}
	}
	return result
}

// All returns the full host-to-tags mapping.
func (t *Tags) All() map[string][]string {
	out := make(map[string][]string, len(t.hostTags))
	for k, v := range t.hostTags {
		out[k] = v
	}
	return out
}
