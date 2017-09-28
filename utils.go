package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Path94/turtleDB"
	"github.com/missionMeteora/toolkit/errors"
)

// common errors returned from the library.
const (
	ErrInvalidToken  = errors.Error("invalid token")
	ErrMissingID     = errors.Error("missing id")
	ErrInvalidLogin  = errors.Error("invalid login")
	ErrNoPassword    = errors.Error("the password is empty")
	ErrNoID          = errors.Error("invalid id")
	ErrUserExists    = errors.Error("user already exists")
	ErrUserNotFound  = errors.Error("user not found")
	ErrBadStatus     = errors.Error("bad status")
	ErrNewUserWithID = errors.Error("a new user can't have an id set")
	ErrPlainPassword = errors.Error("plain password")
)

// marshalUser is used by turtle for marshaling users
func marshalUser(v turtleDB.Value) ([]byte, error) {
	u, ok := v.(User)
	if !ok {
		return nil, unexpectedTypeError(v)
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	if u.ID == "" {
		return nil, ErrNoID
	}

	if u.CreatedTS == 0 {
		// this is a new user, set the Created timestamp.
		u.CreatedTS = time.Now().Unix()
	} else {
		// not a new user, update the LastUpdated timestamp.
		u.LastUpdatedTS = time.Now().Unix()
	}

	return json.Marshal(u)
}

// EditUserTx is a helper func for Auth.EditUser.
func EditUserTx(tx turtleDB.Txn, id string, fn func(u *User) error) (err error) {
	var (
		usersB, _  = tx.Get("users")
		loginsB, _ = tx.Get("logins")

		u User
	)

	if usersB == nil || loginsB == nil {
		// this is a panic because if it happens, something is extremely wrong
		log.Panic("database corruption, can't find bucket")
	}

	if u, err = GetUserByIDTx(tx, id); err != nil {
		return
	}

	// allow changing username
	oldUser := u.Username
	if err = fn(&u); err != nil {
		return
	}

	if err = u.Validate(); err != nil {
		return
	}

	if oldUser != u.Username { // username change
		if oid, _ := GetUserIDTx(tx, u.Username); oid != "" {
			return ErrUserExists
		}

		loginsB.Delete(oldUser)
		loginsB.Put(u.Username, u.ID)
	}

	return usersB.Put(u.ID, u)
}

// GetUserByIDTx is a helper func for Auth.GetUserByID.
func GetUserByIDTx(tx turtleDB.Txn, id string) (usr User, err error) {
	usersB, _ := tx.Get("users")
	if usersB == nil {
		// this is a panic because if it happens, something is extremely wrong
		log.Panic("database corruption, can't find bucket")
	}

	var v turtleDB.Value
	if v, err = usersB.Get(id); err != nil {
		return
	}

	switch v := v.(type) {
	case nil:
		err = ErrUserNotFound
		return
	case User:
		usr = v
		return
	default:
		err = unexpectedTypeError(v)
		return
	}
}

// GetUserByNameTx is a helper func for Auth.GetUserByName.
func GetUserByNameTx(tx turtleDB.Txn, username string) (usr User, err error) {
	var id string
	if id, err = GetUserIDTx(tx, username); err != nil {
		return
	}

	return GetUserByIDTx(tx, id)
}

// GetUserIDTx is a helper func for Auth.GetUserID.
func GetUserIDTx(tx turtleDB.Txn, username string) (string, error) {
	loginsB, _ := tx.Get("logins")
	if loginsB == nil {
		// this is a panic because if it happens, something is extremely wrong
		log.Panic("database corruption, can't find bucket")
	}

	v, err := loginsB.Get(username)
	if err != nil {
		return "", err
	}

	switch v := v.(type) {
	case nil:
		return "", ErrUserNotFound
	case string:
		return v, nil
	default:
		return "", unexpectedTypeError(v)
	}
}

func unexpectedTypeError(v interface{}) error {
	return fmt.Errorf("unexpected type (%T): %#+v", v, v)
}
