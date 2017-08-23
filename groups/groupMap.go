package groups

type groupMap map[string]struct{}

func (g groupMap) Has(group string) (ok bool) {
	_, ok = g[group]
	return
}

func (g groupMap) Set(group string) (ok bool) {
	if g.Has(group) {
		return false
	}

	g[group] = struct{}{}
	return true
}

func (g groupMap) Remove(group string) (ok bool) {
	if ok = g.Has(group); ok {
		delete(g, group)
	}

	return
}

func (g groupMap) Dup() (out groupMap) {
	out = make(groupMap, len(g))
	for group := range g {
		out[group] = struct{}{}
	}

	return
}

func (g groupMap) Slice() (out []string) {
	out = make([]string, 0, len(g))
	for group := range g {
		out = append(out, group)
	}

	return
}
