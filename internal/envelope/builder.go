package envelope

import (
	"fmt"
	"sync/atomic"

	"github.com/user/portwatch/internal/scanner"
)

// Builder constructs Envelopes with a consistent set of default tags,
// labels, and an auto-incrementing scan ID prefix.
type Builder struct {
	prefix     string
	defaultTags   []string
	defaultLabels map[string]string
	counter    atomic.Uint64
}

// NewBuilder returns a Builder whose scan IDs are prefixed with prefix.
func NewBuilder(prefix string) *Builder {
	return &Builder{
		prefix:        prefix,
		defaultLabels: make(map[string]string),
	}
}

// Tag registers a tag that will be attached to every envelope the
// Builder produces.
func (b *Builder) Tag(tag string) *Builder {
	b.defaultTags = append(b.defaultTags, tag)
	return b
}

// Label registers a key/value label that will be attached to every
// envelope the Builder produces.
func (b *Builder) Label(key, value string) *Builder {
	b.defaultLabels[key] = value
	return b
}

// Build wraps the provided host and results in a new Envelope, applying
// all default tags and labels and assigning the next scan ID.
func (b *Builder) Build(host string, results []scanner.Result) *Envelope {
	n := b.counter.Add(1)
	id := fmt.Sprintf("%s-%d", b.prefix, n)

	e := New(host, results).WithScanID(id)

	for _, t := range b.defaultTags {
		e.WithTags(t)
	}
	for k, v := range b.defaultLabels {
		e.WithLabel(k, v)
	}
	return e
}

// Count returns the total number of envelopes produced by this Builder.
func (b *Builder) Count() uint64 {
	return b.counter.Load()
}
