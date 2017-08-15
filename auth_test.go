package auth

import (
	"testing"

	"github.com/Path94/turtleDB"
)

func TestIncID(t *testing.T) {
	d, err := New("/tmp/auth_test")
	if err != nil {
		t.Fatal(err)
	}
	defer d.t.Close()

	if err = d.t.Update(func(tx turtleDB.Txn) error {
		id, err := d.nextID(tx, "auth")
		if err != nil {
			return err
		}
		t.Logf("%s", id)
		id, err = d.nextID(tx, "auth")
		if err != nil {
			return err
		}
		t.Logf("%s", id)
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	if err = d.t.Update(func(tx turtleDB.Txn) error {
		id, err := d.nextID(tx, "auth")
		if err != nil {
			return err
		}
		t.Logf("%s", id)
		id, err = d.nextID(tx, "auth")
		if err != nil {
			return err
		}
		t.Logf("%s", id)
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	t.Log(HashPassword("x"))
}
