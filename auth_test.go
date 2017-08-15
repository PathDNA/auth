package auth

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/Path94/turtleDB"
)

func newTempDB() (*Auth, func(), error) {
	tmpPath, err := ioutil.TempDir("", "auth")
	if err != nil {
		return nil, nil, err
	}
	a, err := New(tmpPath)
	if err != nil {
		return nil, nil, err
	}

	return a, func() {
		a.Close()
		os.RemoveAll(tmpPath)
	}, nil
}

func TestMain(m *testing.M) {
	log.SetFlags(log.Lshortfile)
	os.Exit(m.Run())
}

func TestIncID(t *testing.T) {
	var (
		a, cleanupFn, err = newTempDB()
		id                string
	)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupFn()

	if err = a.t.Update(func(tx turtleDB.Txn) error {
		if id, err = a.nextID(tx, "users"); err != nil {
			return err
		}
		if id, err = a.nextID(tx, "users"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Error(err)
	}

	if err = a.t.Update(func(tx turtleDB.Txn) error {
		if id, err = a.nextID(tx, "users"); err != nil {
			return err
		}
		if id, err = a.nextID(tx, "users"); err != nil {
			return err
		}
		return nil
	}); err != nil {
		t.Error(err)
	}
	if id != "4" {
		t.Errorf("unexpected id: %q", id)
	}
}

func TestCreateUser(t *testing.T) {
	a, cleanupFn, err := newTempDB()
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupFn()

	a.ProfileFn = func() interface{} { return &Profile{} }
	u := &User{
		Status: StatusActive,

		Username: "gbusters",
		Password: "who are you gonna call",

		Profile: &Profile{
			Name:  "Ghost Busters",
			Phone: "1-800-555-2368",
		},
	}

	id, err := a.CreateUser(u, u.Password)
	if isErr(t, err) { // using this because t.Fatal wouldn't run our clean up
		return
	}

	nu, err := a.GetUserByID(u.ID)
	if isErr(t, err) {
		return
	}

	if !reflect.DeepEqual(u, nu) {
		t.Errorf("u != nu\n%#+v\n%#+v", u, nu)
	}
	t.Logf("user (%s): %#+v", id, u)
}

func isErr(t *testing.T, err error) bool {
	t.Helper()
	if err == nil {
		return false
	}
	t.Error(err)
	return true
}
