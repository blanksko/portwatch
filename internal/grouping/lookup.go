package grouping

// GroupsFor returns the names of all groups that contain the given host.
func (r *Registry) GroupsFor(host string) []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var names []string
	for name, hosts := range r.groups {
		for _, h := range hosts {
			if h == host {
				names = append(names, name)
				break
			}
		}
	}
	return names
}

// InGroup reports whether the given host belongs to the named group.
func (r *Registry) InGroup(name, host string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, h := range r.groups[name] {
		if h == host {
			return true
		}
	}
	return false
}
