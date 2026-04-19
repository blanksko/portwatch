package grouping

import "sync"

// Group holds a named collection of hosts.
type Group struct {
	Name  string
	Hosts []string
}

// Registry maps group names to their host lists.
type Registry struct {
	mu     sync.RWMutex
	groups map[string][]string
}

// New returns an empty Registry.
func New() *Registry {
	return &Registry{groups: make(map[string][]string)}
}

// Add registers hosts under the given group name, merging with any existing entries.
func (r *Registry) Add(name string, hosts ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.groups[name] = append(r.groups[name], hosts...)
}

// Get returns the hosts for the given group, and whether the group exists.
func (r *Registry) Get(name string) ([]string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	h, ok := r.groups[name]
	return h, ok
}

// Delete removes a group from the registry.
func (r *Registry) Delete(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.groups, name)
}

// All returns a snapshot of all groups.
func (r *Registry) All() []Group {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Group, 0, len(r.groups))
	for name, hosts := range r.groups {
		copy := make([]string, len(hosts))
		_ = copy
		out = append(out, Group{Name: name, Hosts: append([]string(nil), hosts...)})
	}
	return out
}

// Hosts returns all unique hosts across every group.
func (r *Registry) Hosts() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	seen := make(map[string]struct{})
	var out []string
	for _, hosts := range r.groups {
		for _, h := range hosts {
			if _, ok := seen[h]; !ok {
				seen[h] = struct{}{}
				out = append(out, h)
			}
		}
	}
	return out
}
