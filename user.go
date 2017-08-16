package auth

import (
	"encoding/json"
	"reflect"
	"time"

	"github.com/missionMeteora/toolkit/errors"
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

// User is a system user.
type User struct {
	ID string `json:"id,omitempty"`

	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`

	Status Status `json:"status,omitempty"`

	CreatedTS     int64 `json:"created,omitempty"`
	LastUpdatedTS int64 `json:"lastUpdated,omitempty"`

	Profile interface{} `json:"profile,omitempty"`
}

// UpdatePassword checks if the password is hashed, if not it will hash it and assign the hashed password.
func (u *User) UpdatePassword() error {
	if u.Password == "" {
		return ErrNoPassword
	}

	if IsHashedPass(u.Password) {
		return nil
	}

	p, err := HashPassword(u.Password)
	if err == nil {
		u.Password = p
	}

	return err
}

// Created returns the creation time of the user.
func (u *User) Created() time.Time { return time.Unix(u.CreatedTS, 0) }

// LastUpdated returns the time of the last user update,
// if it was never updated it will return the creation time.
func (u *User) LastUpdated() time.Time {
	if u.LastUpdatedTS == 0 {
		return u.Created()
	}
	return time.Unix(u.LastUpdatedTS, 0)
}

// PasswordsMatch returns true if the current user's hashed password matches the plain-text password.
func (u *User) PasswordsMatch(plainPassword string) bool {
	return CheckPassword(u.Password, plainPassword)
}

// Validate checks if the User struct is valid or not.
func (u *User) Validate() error {
	if u.Password == "" {
		return ErrNoPassword
	}
	if !IsHashedPass(u.Password) {
		return ErrPlainPassword
	}
	if u.Status < StatusActive {
		return ErrBadStatus
	}
	return nil
}

// ErrProfileNotPtr is returned froom UnmarshalUser if profle is not nil and is not a pointer.
const ErrProfileNotPtr = errors.Error("profile must be a pointer")

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
