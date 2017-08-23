package groups

import (
	"github.com/Path94/turtleDB"
	"github.com/missionMeteora/toolkit/errors"
)

const (
	// ErrInvalidType is returned when an invalid type is stored for an id
	ErrInvalidType = errors.Error("invalid type")
	// ErrGroupAlreadySet is returned when a group has been attempted to be set, but already exists for an id
	ErrGroupAlreadySet = errors.Error("cannot set group, group has already been set for this id")
	// ErrGroupNotSet is returned when a group has been attempted to be removed, but doesn't belong for an id
	ErrGroupNotSet = errors.Error("cannot remove group, group does not belong to this id")

	bucketName = "default"
)

var gFuncsMap = turtleDB.NewFuncsMap(marshal, unmarshal)

// New will return a new instance of groups
func New(path string) (gp *Groups, err error) {
	var g Groups
	if g.db, err = turtleDB.New("groups", path, gFuncsMap); err != nil {
		return
	}

	gp = &g
	return
}

// Groups manages user groups
type Groups struct {
	db *turtleDB.Turtle
}

func (g *Groups) get(txn turtleDB.Txn, id string) (gm groupMap, err error) {
	var (
		bkt turtleDB.Bucket
		val turtleDB.Value
		ok  bool
	)

	if bkt, err = txn.Get(bucketName); err != nil {
		return
	}

	if val, err = bkt.Get(id); err != nil {
		return
	}

	if gm, ok = val.(groupMap); !ok {
		err = ErrInvalidType
		return
	}

	return
}
func (g *Groups) put(txn turtleDB.Txn, id string, gm groupMap) (err error) {
	var bkt turtleDB.Bucket

	if bkt, err = txn.Create(bucketName); err != nil {
		return
	}

	if err = bkt.Put(id, gm); err != nil {
		return
	}

	return
}

// Get will get the group slice by id
func (g *Groups) Get(id string) (gs []string, err error) {
	var gm groupMap
	err = g.db.Read(func(txn turtleDB.Txn) (err error) {
		if gm, err = g.get(txn, id); err != nil {
			return
		}

		gs = gm.Slice()
		return
	})

	return
}

// Has will confirm if an id has a given group
func (g *Groups) Has(id, group string) (has bool) {
	var gm groupMap
	g.db.Read(func(txn turtleDB.Txn) (err error) {
		if gm, err = g.get(txn, id); err != nil {
			return
		}

		has = gm.Has(group)
		return
	})

	return
}

// Set will set a group to a given id
func (g *Groups) Set(id string, groups ...string) (gs []string, err error) {
	var (
		gm groupMap
	)

	err = g.db.Update(func(txn turtleDB.Txn) (err error) {
		if gm, err = g.get(txn, id); err != nil {
			gm = make(groupMap)
			err = nil
		}

		for _, group := range groups {
			if !gm.Set(group) {
				err = ErrGroupAlreadySet
				return
			}

			if err = g.put(txn, id, gm.Dup()); err != nil {
				return
			}
		}

		return
	})

	return gm.Slice(), err
}

// Remove will remove a group from a given id
func (g *Groups) Remove(id, group string) (gs []string, err error) {
	var gm groupMap
	err = g.db.Update(func(txn turtleDB.Txn) (err error) {
		if gm, err = g.get(txn, id); err != nil {
			gm = make(groupMap)
			err = nil
		}

		if !gm.Remove(group) {
			return ErrGroupNotSet
		}

		return g.put(txn, id, gm.Dup())
	})

	gs = gm.Slice()
	return
}

// Close will close groups
func (g *Groups) Close() (err error) {
	return g.db.Close()
}
