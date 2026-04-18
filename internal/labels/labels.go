// Package labels provides host labelling for portwatch.
package labels

import "sync"

// Labels holds a map of host → free-form label string.
type Labels struct {
	mu   sync.RWMutex
	data map[string]string
}

// New returns an empty Labels store pre-populated with the provided map.
// A nil map is accepted and treated as empty.
func New(initial map[string]string) *Labels {
	l := &Labels{data: make(map[string]string, len(initial))}
	for k, v := range initial {
		l.data[k] = v
	}
	return l
}

// Set assigns a label to a host, replacing any existing value.
func (l *Labels) Set(host, label string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.data[host] = label
}

// Get returns the label for host and whether it was found.
func (l *Labels) Get(host string) (string, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()
	v, ok := l.data[host]
	return v, ok
}

// Delete removes the label for host.
func (l *Labels) Delete(host string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.data, host)
}

// All returns a shallow copy of all host→label mappings.
func (l *Labels) All() map[string]string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make(map[string]string, len(l.data))
	for k, v := range l.data {
		out[k] = v
	}
	return out
}
