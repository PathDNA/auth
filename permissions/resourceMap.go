package permissions

type resourceMap map[string]uint8

// Get will get the permissions available to a given group
func (r resourceMap) Get(group string) (permissions uint8, ok bool) {
	permissions, ok = r[group]
	return
}

// Set will set the permissions available to a given group
func (r resourceMap) Set(group string, permissions uint8) (ok bool) {
	var currentPermissions uint8
	if currentPermissions, _ = r.Get(group); currentPermissions == permissions {
		return false
	}

	r[group] = permissions
	return true
}

func (r resourceMap) Remove(group string) (ok bool) {
	if _, ok = r.Get(group); ok {
		delete(r, group)
	}

	return
}

func (r resourceMap) Dup() (out resourceMap) {
	out = make(resourceMap, len(r))
	for group, permissions := range r {
		out[group] = permissions
	}

	return
}
