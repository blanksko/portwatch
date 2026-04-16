package notifier

import (
	"fmt"
	"io"

	"github.com/user/portwatch/internal/snapshot"
)

// Multi fans out notifications to multiple Notifier instances.
type Multi struct {
	notifiers []*Notifier
}

// NewMulti creates a Multi notifier writing to each of the supplied writers.
// It returns an error if any writer is nil.
func NewMulti(writers ...io.Writer) (*Multi, error) {
	ns := make([]*Notifier, 0, len(writers))
	for i, w := range writers {
		n, err := New(w)
		if err != nil {
			return nil, fmt.Errorf("notifier: writer at index %d is nil", i)
		}
		ns = append(ns, n)
	}
	return &Multi{notifiers: ns}, nil
}

// Notify calls Notify on every contained Notifier, collecting all errors.
func (m *Multi) Notify(host string, d snapshot.DiffResult) error {
	var first error
	for _, n := range m.notifiers {
		if err := n.Notify(host, d); err != nil && first == nil {
			first = err
		}
	}
	return first
}
