package permissions

import (
	"encoding/json"

	"github.com/PathDNA/turtleDB"
)

type resource map[string]Action

// Get will get the actions available to a given group
func (r resource) Get(group string) (actions Action, ok bool) {
	actions, ok = r[group]
	return
}

// Set will set the actions available to a given group
func (r resource) Set(group string, actions Action) (ok bool) {
	var currentActions Action
	currentActions, _ = r.Get(group)
	if currentActions|actions == currentActions {
		return false
	}

	r[group] = actions
	return true
}

func (r resource) canRead(group string) bool {
	act := r[group]
	return act&ActionRead != 0
}

func (r resource) canWrite(group string) bool {
	act := r[group]
	return act&ActionWrite != 0
}

func (r resource) canDelete(group string) bool {
	act := r[group]
	return act&ActionDelete != 0
}

func (r resource) Remove(group string) (ok bool) {
	if _, ok = r.Get(group); ok {
		delete(r, group)
	}

	return
}

func (r resource) Dup() (out resource) {
	out = make(resource, len(r))
	for group, actions := range r {
		out[group] = actions
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
