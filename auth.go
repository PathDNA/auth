package auth

import (
	"encoding/json"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/itsmontoya/middleware"

	"github.com/PathDNA/turtleDB"
)

var (
	buckets = &[...]string{"users", "logins", "tokens", "index"}

	one = big.NewInt(1)
)

// Auth is a generic user authentication helper.
type Auth struct {
	t *turtleDB.Turtle

	//ProfileFn is used on loading users from the database to fill in the User.Profile field.}
	profileFn atomic.Value
}

// New returns a new Auth db at the specificed path.
func New(path string) (*Auth, error) {
	return NewEncrypted(path, nil, nil)
}

// NewEncrypted returns a new Auth db that is encrypted with the specified key/iv.
// if key is nil, it returns a non-encrypted store.
func NewEncrypted(path string, key, iv []byte) (*Auth, error) {
	var (
		a       Auth
		funcMap = turtleDB.NewFuncsMap(turtleDB.MarshalJSON, turtleDB.UnmarshalJSON)
		err     error
	)

	funcMap.Put("users", marshalUser, a.unmarshalUser)

	if key != nil {
		a.t, err = turtleDB.New("auth", path, funcMap, middleware.NewCryptyMW(key, iv))
	} else {
		a.t, err = turtleDB.New("auth", path, funcMap)
	}

	if err != nil {
		return nil, err
	}

	if err = a.t.Update(func(tx turtleDB.Txn) error {
		for _, b := range buckets {
			if _, err = tx.Create(b); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &a, nil
}

// NewProfileFn is used on loading users from the database to fill in the User.Profile field.
// it is 100% optional
func (a *Auth) NewProfileFn(fn func() interface{}) {
	a.profileFn.Store(fn)
}

func (a *Auth) getProfileFn() func() interface{} {
	fn, _ := a.profileFn.Load().(func() interface{})
	return fn
}

// CreateUser will add the passed user to the database and hash the given password.
// the passed user will be modified with the hashed password and the new ID.
func (a *Auth) CreateUser(username, password string) (id string, err error) {
	var u User
	// hash outside the db lock
	if u.Password, err = HashPassword(password); err != nil {
		return
	}

	u.Status = StatusInactive
	u.Username = username
	u.CreatedTS = time.Now().Unix()
	u.LastUpdatedTS = u.CreatedTS

	if err = u.Validate(); err != nil {
		return
	}

	if err = a.t.Update(func(tx turtleDB.Txn) error {
		var (
			loginsB, _ = tx.Get("logins")
			usersB, _  = tx.Get("users")
		)

		if id, _ := GetUserIDTx(tx, u.Username); id != "" {
			return ErrUserExists
		}

		if u.ID, err = a.nextID(tx, "users"); err != nil {
			return err
		}

		if err = usersB.Put(u.ID, u); err != nil {
			return err
		}

		return loginsB.Put(u.Username, u.ID)
	}); err != nil {
		return
	}

	id = u.ID
	return
}

// EditUserByID edits a user by their ID, returning an error will cancel the edit.
func (a *Auth) EditUserByID(id string, fn func(u *User) error) error {
	return a.t.Update(func(tx turtleDB.Txn) error {
		return EditUserTx(tx, id, fn)
	})
}

// EditUserByName edits a user by their username, returning an error will cancel the edit.
func (a *Auth) EditUserByName(username string, fn func(u *User) error) error {
	return a.t.Update(func(tx turtleDB.Txn) error {
		id, err := GetUserIDTx(tx, username)
		if err != nil {
			return err
		}
		return EditUserTx(tx, id, fn)
	})
}

// GetUserByID returns a User by their ID.
func (a *Auth) GetUserByID(id string) (u User, err error) {
	err = a.t.Read(func(tx turtleDB.Txn) error {
		u, err = GetUserByIDTx(tx, id)
		return err
	})
	return
}

// GetUserByName returns a User by their UserName.
func (a *Auth) GetUserByName(username string) (u User, err error) {
	err = a.t.Read(func(tx turtleDB.Txn) error {
		u, err = GetUserByNameTx(tx, username)
		return err
	})
	return
}

// ForEach will iterate through each of the users
func (a *Auth) ForEach(fn func(User) error) (err error) {
	return a.t.Read(func(txn turtleDB.Txn) (err error) {
		var bkt turtleDB.Bucket
		if bkt, err = txn.Get("users"); err != nil {
			return
		}

		return bkt.ForEach(func(key string, val turtleDB.Value) (err error) {
			var (
				u  User
				ok bool
			)

			if u, ok = val.(User); !ok {
				return turtleDB.ErrInvalidType
			}

			return fn(u)
		})
	})
}

// Close closes the underlying database.
func (a *Auth) Close() error {
	return a.t.Close()
}

// unmarshalUser is a helper for turtleDB.
func (a *Auth) unmarshalUser(p []byte) (turtleDB.Value, error) {
	var u User

	if pfn := a.getProfileFn(); pfn != nil {
		u.Profile = pfn()
	}

	if err := json.Unmarshal(p, &u); err != nil {
		return nil, err
	}

	return u, nil
}

func (a *Auth) nextID(tx turtleDB.Txn, bucket string) (string, error) {
	b, err := tx.Get("index")
	if err != nil {
		return "", err
	}

	v, err := b.Get(bucket)
	if err != nil && err != turtleDB.ErrKeyDoesNotExist {
		return "", err
	}

	n := big.NewInt(0)
	switch v := v.(type) {
	case nil:
	case string:
		if _, ok := n.SetString(v, 10); !ok {
			return "", unexpectedTypeError(v)
		}

	default:
		return "", unexpectedTypeError(v)
	}

	id := n.Add(n, one).String()
	if err = b.Put(bucket, id); err != nil {
		return "", err
	}

	return id, nil
}
