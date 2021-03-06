package permissions

import (
	"github.com/PathDNA/turtleDB"
	"github.com/missionMeteora/toolkit/errors"
)

const (
	// ErrInvalidActions is returned when an invalid permissions value is attempted to be set
	ErrInvalidActions = errors.Error("invalid permissions, please see constant block for reference")
	// ErrPermissionsUnchanged is returned when matching permissions are set for a resource
	ErrPermissionsUnchanged = errors.Error("permissions match, unchanged")
)

// Action represents an action type
type Action uint8

// Can will return if an action can peform an action request
func (a Action) Can(ar Action) (can bool) {
	return a&ar != 0
}

const (
	// ActionNone represents a zero value, no action
	ActionNone Action = 1 << iota
	// ActionRead represents a reading action
	ActionRead
	// ActionWrite represents a writing action
	ActionWrite
	// ActionDelete represents a deleting action
	ActionDelete
)

const (
	resourceBkt = "resource"
	groupsBkt   = "groups"
)

// New will return a new instance of Permissions
func New(dir string) (pp *Permissions, err error) {
	var p Permissions
	if err = p.initDB(dir); err != nil {
		return
	}

	pp = &p
	return
}

// Permissions manages permissions
type Permissions struct {
	db *turtleDB.Turtle
}

func (p *Permissions) initDB(dir string) (err error) {
	fm := turtleDB.NewFuncsMap(marshalResource, unmarshalResource)
	fm.Put(groupsBkt, marshalGroups, unmarshalGroups)

	if p.db, err = turtleDB.New("permissions", dir, fm); err != nil {
		return
	}

	return p.db.Update(func(txn turtleDB.Txn) (err error) {
		if _, err = txn.Create(resourceBkt); err != nil {
			return
		}

		_, err = txn.Create(groupsBkt)
		return
	})
}

func (p *Permissions) getResource(txn turtleDB.Txn, id string) (r resource, err error) {
	var (
		bkt turtleDB.Bucket
		val turtleDB.Value
		ok  bool
	)

	if bkt, err = txn.Get(resourceBkt); err != nil {
		return
	}

	if val, err = bkt.Get(id); err != nil {
		return
	}

	if r, ok = val.(resource); !ok {
		err = turtleDB.ErrInvalidType
		return
	}

	return
}

func (p *Permissions) getGroups(txn turtleDB.Txn, uuid string) (g groups, err error) {
	var (
		bkt turtleDB.Bucket
		val turtleDB.Value
		ok  bool
	)

	if bkt, err = txn.Get(groupsBkt); err != nil {
		return
	}

	if val, err = bkt.Get(uuid); err != nil {
		return
	}

	if g, ok = val.(groups); !ok {
		err = turtleDB.ErrInvalidType
		return
	}

	return
}

func (p *Permissions) putResource(txn turtleDB.Txn, id string, r resource) (err error) {
	var bkt turtleDB.Bucket
	if bkt, err = txn.Get(resourceBkt); err != nil {
		return
	}

	if err = bkt.Put(id, r); err != nil {
		return
	}

	return
}

func (p *Permissions) putGroups(txn turtleDB.Txn, id string, g groups) (err error) {
	var bkt turtleDB.Bucket
	if bkt, err = txn.Get(groupsBkt); err != nil {
		return
	}

	if err = bkt.Put(id, g); err != nil {
		return
	}

	return
}

// Get will get the permissions for a given group for a resource id
func (p *Permissions) Get(id, group string) (actions Action) {
	var r resource
	p.db.Read(func(txn turtleDB.Txn) (err error) {
		if r, err = p.getResource(txn, id); err != nil {
			if err == turtleDB.ErrKeyDoesNotExist {
				err = nil
			}
			return
		}

		actions, _ = r.Get(group)
		return
	})

	return
}

// SetPermissions will set the permissions for a given group for a resource id
func (p *Permissions) SetPermissions(id, group string, actions Action) (err error) {
	var r resource
	if !isValidActions(actions) {
		return ErrInvalidActions
	}

	return p.db.Update(func(txn turtleDB.Txn) (err error) {
		if r, err = p.getResource(txn, id); err != nil {
			if err != turtleDB.ErrKeyDoesNotExist {
				return
			}

			r = make(resource)
			err = nil
		} else {
			r = r.Dup()
		}

		if !r.Set(group, actions) {
			return ErrPermissionsUnchanged
		}

		err = p.putResource(txn, id, r)
		return
	})
}

// AddGroup will add a group to a uuid
func (p *Permissions) AddGroup(uuid string, grouplist ...string) (err error) {
	var g groups
	return p.db.Update(func(txn turtleDB.Txn) (err error) {
		if g, err = p.getGroups(txn, uuid); err != nil {
			if err != turtleDB.ErrKeyDoesNotExist {
				return
			}

			g = make(groups)
			err = nil
		} else {
			g = g.Dup()
		}

		updated := false
		for _, group := range grouplist {
			if g.Set(group) {
				updated = true
			}
		}

		if !updated {
			return ErrPermissionsUnchanged
		}

		return p.putGroups(txn, uuid, g)
	})
}

// RemoveGroup will remove a group to a uuid
func (p *Permissions) RemoveGroup(uuid string, grouplist ...string) (err error) {
	var g groups
	return p.db.Update(func(txn turtleDB.Txn) (err error) {
		if g, err = p.getGroups(txn, uuid); err != nil {
			return ErrPermissionsUnchanged
		}

		g = g.Dup()
		updated := false

		for _, group := range grouplist {
			if g.Remove(group) {
				updated = true
			}
		}

		if !updated {
			return ErrPermissionsUnchanged
		}

		return p.putGroups(txn, uuid, g)
	})
}

// Can will return if a user (UUID) can perform a given action on a provided resource id
func (p *Permissions) Can(uuid, id string, action Action) (can bool) {
	var (
		g   groups
		r   resource
		err error
	)

	if err = p.db.Read(func(txn turtleDB.Txn) (err error) {
		if g, err = p.getGroups(txn, uuid); err != nil {
			return
		}

		if r, err = p.getResource(txn, id); err != nil {
			return
		}

		switch action {
		case ActionNone:
			can = true

		case ActionRead:
			can = g.ForEach(r.canRead)

		case ActionWrite:
			can = g.ForEach(r.canWrite)

		case ActionDelete:
			can = g.ForEach(r.canDelete)
		}

		return
	}); err != nil {
		return
	}

	return
}

// Has will return whether or not an ID has a particular group associated with it
func (p *Permissions) Has(id, group string) (ok bool) {
	var g groups
	p.db.Read(func(txn turtleDB.Txn) (err error) {
		if g, err = p.getGroups(txn, id); err != nil {
			return
		}

		ok = g.Has(group)
		return
	})

	return
}

// Groups will return a slice of the groups a user belongs to
func (p *Permissions) Groups(uuid string) (gs []string, err error) {
	var g groups

	err = p.db.Read(func(txn turtleDB.Txn) (err error) {
		if g, err = p.getGroups(txn, uuid); err != nil {
			return
		}

		gs = g.Slice()
		return
	})

	return
}

// Close will close permissions
func (p *Permissions) Close() (err error) {
	return p.db.Close()
}
