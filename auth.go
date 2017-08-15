package auth

import (
	"math/big"

	"github.com/Path94/turtleDB"
)

var (
	buckets = &[...]string{"auth", "login", "tokens"}
	one     = big.NewInt(1)
)

type Auth struct {
	t *turtleDB.Turtle
}

func New(path string) (*Auth, error) {
	t, err := turtleDB.New("auth", path, nil)
	if err != nil {
		return nil, err
	}

	// if err = t.Update(func(tx turtleDB.Txn) error {
	// 	for _, b := range buckets {
	// 		if _, err := tx.Create(b); err != nil {
	// 			return err
	// 		}
	// 	}
	// 	return nil
	// }); err != nil {
	// 	return nil, err
	// }

	return &Auth{t: t}, nil
}

func (a *Auth) nextID(tx turtleDB.Txn, bucket string) (*big.Int, error) {
	const idCounterKey = ":id:"

	b, err := tx.Get(bucket)
	if err != nil {
		return nil, err
	}

	v, err := b.Get(idCounterKey)
	if err != nil && err != turtleDB.ErrKeyDoesNotExist {
		return nil, err
	}

	n, _ := v.(*big.Int)

	if n == nil {
		n = big.NewInt(0)
	}

	if err = b.Put(idCounterKey, n.Add(n, one)); err != nil {
		return nil, err
	}

	return n, nil
}

func (a *Auth) CreateUser(u *User, password string) (err error) {
	if u.ID != nil {
		return ErrNewUserWithID
	}

	// hash outside the db lock
	if u.Password, err = HashPassword(password); err != nil {
		return
	}

	return a.t.Update(func(tx turtleDB.Txn) error {
		var (
			loginsB, _ = tx.Get("logins")
			authB, _   = tx.Get("auth")
		)

		if _, err := loginsB.Get(u.Username); err != turtleDB.ErrKeyDoesNotExist {
			return ErrUserExists
		}

		if u.ID, err = a.nextID(tx, "auth"); err != nil {
			return err
		}

		id := u.ID.String()

		if err = authB.Put(id, u); err != nil {
			return err
		}

		return loginsB.Put(u.Username, id)
	})
}
