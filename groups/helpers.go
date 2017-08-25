package groups

import (
	"encoding/json"

	"github.com/Path94/turtleDB"
)

func marshal(val turtleDB.Value) (b []byte, err error) {
	var (
		gm groupMap
		ok bool
	)

	if gm, ok = val.(groupMap); !ok {
		err = ErrInvalidType
		return
	}

	return json.Marshal(gm.Slice())
}

func unmarshal(b []byte) (val turtleDB.Value, err error) {
	var (
		gs []string
		gm groupMap
	)

	if err = json.Unmarshal(b, &gs); err != nil {
		return
	}

	gm = make(groupMap, len(gs))
	for _, group := range gs {
		gm.Set(group)
	}

	val = gm
	return
}
