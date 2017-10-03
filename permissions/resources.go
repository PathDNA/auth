package permissions

import (
	"encoding/json"

	"github.com/PathDNA/turtleDB"
)

type resource map[string]uint8

// Get will get the permissions available to a given group
func (r resource) Get(group string) (permissions uint8, ok bool) {
	permissions, ok = r[group]
	return
}

// Set will set the permissions available to a given group
func (r resource) Set(group string, permissions uint8) (ok bool) {
	var currentPermissions uint8
	if currentPermissions, _ = r.Get(group); currentPermissions == permissions {
		return false
	}

	r[group] = permissions
	return true
}

func (r resource) canRead(group string) bool {
	perm := r[group]
	return perm == PermissionRead || perm == PermissionReadWrite
}

func (r resource) canWrite(group string) bool {
	perm := r[group]
	return perm == PermissionWrite || perm == PermissionReadWrite
}

func (r resource) Remove(group string) (ok bool) {
	if _, ok = r.Get(group); ok {
		delete(r, group)
	}

	return
}

func (r resource) Dup() (out resource) {
	out = make(resource, len(r))
	for group, permissions := range r {
		out[group] = permissions
	}

	return
}

func marshalResource(val turtleDB.Value) (b []byte, err error) {
	var (
		r  resource
		ok bool
	)

	if r, ok = val.(resource); !ok {
		err = turtleDB.ErrInvalidType
		return
	}

	return json.Marshal(r)
}

func unmarshalResource(b []byte) (val turtleDB.Value, err error) {
	var (
		r resource
	)

	if err = json.Unmarshal(b, &r); err != nil {
		return
	}

	val = r
	return
}
