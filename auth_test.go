package auth

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/PathDNA/turtleDB"
)

func newTempDB(enc bool) (a *Auth, cleanup func(), err error) {
	tmpPath, err := ioutil.TempDir("", "auth")
	if err != nil {
		return nil, nil, err
	}
	if enc {
		var (
			iv  = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
			key = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
		)
		a, err = NewEncrypted(tmpPath, key, iv)
	} else {
		a, err = New(tmpPath)
	}
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
		a, cleanupFn, err = newTempDB(false)
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
	t.Run("Plain", func(t *testing.T) {
		testCreateUser(t, false)
	})

	t.Run("Encrypted", func(t *testing.T) {
		testCreateUser(t, true)
	})
}

func testCreateUser(t *testing.T, enc bool) {
	a, cleanupFn, err := newTempDB(enc)
	if err != nil {
		t.Fatal(err)
	}
	defer cleanupFn()

	a.NewProfileFn(func() interface{} { return &Profile{} })
	u := &User{
		Status: StatusActive,

		Username: "gbusters",
		Password: "who are you gonna call",

		Profile: Profile{
			Name:  "Ghost Busters",
			Phone: "1-800-555-2368",
		},
	}

	id, err := a.CreateUser("gbusters", u.Password)
	if isErr(t, err) { // using this because t.Fatal wouldn't run our clean up
		return
	}

	a.EditUserByID(id, func(usr *User) (err error) {
		usr.Status = StatusActive
		usr.Profile = u.Profile
		return
	})

	nu, err := a.GetUserByID(id)
	if isErr(t, err) {
		return
	}

	if u.Status != nu.Status {
		t.Fatalf("status does not match: %v / %v", u.Status, nu.Status)
	}

	if u.Username != nu.Username {
		t.Fatalf("username does not match: %v / %v", u.Username, nu.Username)
	}

	if !nu.PasswordsMatch(u.Password) {
		t.Fatalf("password does not match: %v", nu.Password)
	}

	if u.Profile != nu.Profile {
		t.Fatalf("profile does not match: %v / %v", u.Profile, nu.Profile)
	}

	var ucnt int
	if err = a.ForEach(func(u User) (err error) {
		ucnt++
		return
	}); err != nil {
		t.Fatal(err)
	}

	if ucnt != 1 {
		t.Fatalf("invalid user count, expected %v and received %v", 1, ucnt)
	}
}

func isErr(t *testing.T, err error) bool {
	t.Helper()
	if err == nil {
		return false
	}
	t.Error(err)
	return true
}
