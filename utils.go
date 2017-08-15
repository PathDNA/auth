package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Path94/turtleDB"
)

func marshalUser(v turtleDB.Value) ([]byte, error) {
	u, ok := v.(*User)
	if !ok {
		return nil, unexpectedTypeError(v)
	}
	if u.Password == "" {
		return nil, ErrNoPassword
	}

	if u.ID == "" {
		return nil, ErrNoID
	}

	if u.Status < StatusActive {
		return nil, ErrBadStatus
	}

	if u.CreatedTS == 0 {
		u.CreatedTS = time.Now().Unix()
	} else {
		u.LastUpdatedTS = time.Now().Unix()
	}

	return json.Marshal(u)
}

func EditUserTx(tx turtleDB.Txn, id string, fn func(u *User) error) (err error) {
	var (
		usersB, _  = tx.Get("users")
		loginsB, _ = tx.Get("logins")
		u          *User
	)
	if usersB == nil || loginsB == nil {
		// this is a panic because if it happens, something is extremely wrong
		log.Panic("database corruption, can't find bucket")
	}
	if u, err = GetUserByIDTx(tx, id); err != nil {
		return
	}

	// allow changing username
	oldUser, oldPass := u.Username, u.Password
	if err = fn(u); err != nil {
		return err
	}

	if oldUser != u.Username { // username change
		if _, err := GetUserIDTx(tx, u.Username); err != nil {
			return ErrUserExists
		}
		loginsB.Delete(oldUser)
		loginsB.Put(u.Username, u.ID)
	}

	if oldPass != u.Password {
		if !IsHashedPass(u.Password) {
			if u.Password, err = HashPassword(u.Password); err != nil {
				return
			}
		}
	}
	return usersB.Put(u.ID, u)
}

func GetUserByIDTx(tx turtleDB.Txn, id string) (*User, error) {
	usersB, _ := tx.Get("users")
	if usersB == nil {
		// this is a panic because if it happens, something is extremely wrong
		log.Panic("database corruption, can't find bucket")
	}
	v, err := usersB.Get(id)
	if err != nil {
		return nil, err
	}

	switch v := v.(type) {
	case nil:
		return nil, ErrUserNotFound
	case *User:
		return v, nil
	default:
		return nil, unexpectedTypeError(v)
	}
}

func GetUserByNameTx(tx turtleDB.Txn, username string) (*User, error) {
	id, err := GetUserIDTx(tx, username)
	if err != nil {
		return nil, err
	}
	return GetUserByIDTx(tx, id)
}

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
