package permissions

import (
	"encoding/json"

	"github.com/Path94/turtleDB"
)

type groups map[string]struct{}

func (g groups) Has(group string) (ok bool) {
	_, ok = g[group]
	return
}

func (g groups) Set(group string) (ok bool) {
	if g.Has(group) {
		return false
	}

	g[group] = struct{}{}
	return true
}

func (g groups) Remove(group string) (ok bool) {
	if ok = g.Has(group); ok {
		delete(g, group)
	}

	return
}

func (g groups) Dup() (out groups) {
	out = make(groups, len(g))
	for group := range g {
		out[group] = struct{}{}
	}

	return
}

func (g groups) Slice() (out []string) {
	out = make([]string, 0, len(g))
	for group := range g {
		out = append(out, group)
	}

	return
}

func (g groups) ForEach(fn func(group string) (end bool)) (ended bool) {
	for group := range g {
		if fn(group) {
			return true
		}
	}

	return
}

func marshalGroups(val turtleDB.Value) (b []byte, err error) {
	var (
		g  groups
		ok bool
	)

	if g, ok = val.(groups); !ok {
		err = turtleDB.ErrInvalidType
		return
	}

	return json.Marshal(g)
}

func unmarshalGroups(b []byte) (val turtleDB.Value, err error) {
	var g groups
	if err = json.Unmarshal(b, &g); err != nil {
		return
	}

	val = g
	return
}
