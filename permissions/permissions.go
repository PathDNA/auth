package permissions

import (
	"github.com/Path94/turtleDB"
	"github.com/missionMeteora/toolkit/errors"
)

const (
	// ErrInvalidPermissions is returned when an invalid permissions value is attempted to be set
	ErrInvalidPermissions = errors.Error("invalid permissions, please see constant block for reference")
	// ErrPermissionsUnchanged is returned when matching permissions are set for a resource
	ErrPermissionsUnchanged = errors.Error("permissions match, unchanged")
)

const (
	// PermissionNone represents a zero value, no permissions
	PermissionNone uint8 = iota
	// PermissionRead represents read-only permissions
	PermissionRead
	// PermissionWrite represents write-only permissions
	PermissionWrite
	// PermissionReadWrite represents read/write permissions
	PermissionReadWrite
)

const (
	// ActionNone represents a zero value, no action
	ActionNone uint8 = iota
	// ActionRead represents a reading action
	ActionRead
	// ActionWrite represents a writing action
	ActionWrite
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
func (p *Permissions) Get(id, group string) (permissions uint8) {
	var r resource
	p.db.Read(func(txn turtleDB.Txn) (err error) {
		if r, err = p.getResource(txn, id); err != nil {
			if err == turtleDB.ErrKeyDoesNotExist {
				err = nil
			}
			return
		}

		permissions, _ = r.Get(group)
		return
	})

	return
}

// SetPermissions will set the permissions for a given group for a resource id
func (p *Permissions) SetPermissions(id, group string, permissions uint8) (err error) {
	var r resource
	if !isValidPermissions(permissions) {
		return ErrInvalidPermissions
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

		if !r.Set(group, permissions) {
			return ErrPermissionsUnchanged
		}

		err = p.putResource(txn, id, r)
		return
	})
}

// AddGroup will add a group to a uuid
func (p *Permissions) AddGroup(uuid string, group string) (err error) {
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

		if !g.Set(group) {
			return ErrPermissionsUnchanged
		}

		err = p.putGroups(txn, uuid, g)
		return
	})
}

// Can will return if a user (UUID) can perform a given action on a provided resource id
func (p *Permissions) Can(uuid, id string, action uint8) (can bool) {
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
		}

		return
	}); err != nil {
		return
	}

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
