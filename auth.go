package auth

import (
	"encoding/json"
	"math/big"

	"github.com/Path94/turtleDB"
)

var (
	buckets = &[...]string{"users", "logins", "tokens", "index"}

	one = big.NewInt(1)
)

type Auth struct {
	t *turtleDB.Turtle

	//ProfileFn is used on loading users from the database to fill in the User.Profile field.}
	ProfileFn func() interface{}
}

func New(path string) (*Auth, error) {
	var a Auth
	funcMap := turtleDB.NewFuncsMap(turtleDB.MarshalJSON, turtleDB.UnmarshalJSON)
	funcMap.Put("users", marshalUser, a.unmarshalUser)

	t, err := turtleDB.New("auth", path, funcMap)
	if err != nil {
		return nil, err
	}

	if err = t.Update(func(tx turtleDB.Txn) error {
		for _, b := range buckets {
			if _, err := tx.Create(b); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	a.t = t

	return &a, nil
}

func (a *Auth) Close() error {
	return a.t.Close()
}

// CreateUser will add the passed user to the database and hash the given password.
// the passed user will be modified with the hashed password and the new ID.
func (a *Auth) CreateUser(u *User, password string) (id string, err error) {
	if u.ID != "" {
		return "", ErrNewUserWithID
	}

	// hash outside the db lock
	if u.Password, err = HashPassword(password); err != nil {
		return
	}

	return u.ID, a.t.Update(func(tx turtleDB.Txn) error {
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
	})
}

func (a *Auth) EditUserByID(id string, fn func(u *User) error) error {
	return a.t.Update(func(tx turtleDB.Txn) error {
		return EditUserTx(tx, id, fn)
	})
}

func (a *Auth) EditUserByName(username string, fn func(u *User) error) error {
	return a.t.Update(func(tx turtleDB.Txn) error {
		id, err := GetUserIDTx(tx, username)
		if err != nil {
			return err
		}
		return EditUserTx(tx, id, fn)
	})
}

func (a *Auth) GetUserByID(id string) (u *User, err error) {
	err = a.t.Read(func(tx turtleDB.Txn) error {
		u, err = GetUserByIDTx(tx, id)
		return err
	})
	return
}

func (a *Auth) GetUserByName(username string) (u *User, err error) {
	err = a.t.Read(func(tx turtleDB.Txn) error {
		u, err = GetUserByNameTx(tx, username)
		return err
	})
	return
}

func (a *Auth) unmarshalUser(p []byte) (turtleDB.Value, error) {
	var u User

	if a.ProfileFn != nil {
		u.Profile = a.ProfileFn()
	}

	if err := json.Unmarshal(p, &u); err != nil {
		return nil, err
	}

	return &u, nil
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
