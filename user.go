package auth

import (
	"encoding/json"
	"math/big"
	"reflect"
	"time"

	"github.com/missionMeteora/toolkit/errors"
)

const (
	ErrProfileNotPtr = errors.Error("profile must be a pointer")
	ErrNoPassword    = errors.Error("the password is empty")
	ErrNoID          = errors.Error("invalid id")
	ErrUserExists = errors.Error("user already exists")
	ErrBadStatus     = errors.Error("bad status")
	ErrNewUserWithID = errors.Error("a new user can't have an id set")
)

// Status represents different user statuses
type Status int8

// Status values.
const (
	_ Status = iota
	StatusActive
	StatusInactive
	StatusBanned
)

// User is a system user
type User struct {
	ID *big.Int `json:"id,omitempty"`

	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`

	Status Status `json:"status,omitempty"`

	CreatedTS     int64 `json:"created,omitempty"`
	LastUpdatedTS int64 `json:"lastUpdated,omitempty"`

	Profile interface{} `json:"profile,omitempty"`
}

// Created returns the creation time of the user.
func (u *User) Created() time.Time { return time.Unix(u.CreatedTS, 0) }

// LastUpdated returns the time of the last user update.
func (u *User) LastUpdated() time.Time { return time.Unix(u.LastUpdatedTS, 0) }

type user_ User // used for marshaling

// MarshalJSON implements json.Marshaler
func (u *User) MarshalJSON() ([]byte, error) {
	if u.Password == "" {
		return nil, ErrNoPassword
	}

	if u.ID == nil {
		return nil, ErrNoID
	}

	if u.Status < StatusActive {
		return nil, ErrBadStatus
	}

	uu := (user_)(*u)

	if uu.CreatedTS == 0 {
		uu.CreatedTS = time.Now().Unix()
	} else {
		uu.LastUpdatedTS = time.Now().Unix()
	}

	return json.Marshal(&uu)
}

// UnmarshalUser attempts to unmarshal json with the optional Profile field and returns the *User.
// if profile is NOT nil, it must be a pointer.
func UnmarshalUser(b []byte, profile interface{}) (*User, error) {
	var u User
	if profile != nil {
		if reflect.TypeOf(profile).Kind() != reflect.Ptr {
			// this should probably panic because that's 100% a programming error
			return nil, ErrProfileNotPtr
		}
		u.Profile = profile
	}

	if err := json.Unmarshal(b, &u); err != nil {
		return nil, err
	}

	return &u, nil
}
