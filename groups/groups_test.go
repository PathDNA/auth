package groups

import (
	"os"
	"testing"
)

func TestGroups(t *testing.T) {
	var (
		g   *Groups
		gs  []string
		err error
	)

	if g, err = New("./_testdata"); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("./_testdata")

	if gs, err = g.Set("0", "users"); err != nil {
		t.Fatalf("Error encountered while setting: %v", err)
	}

	if len(gs) != 1 {
		t.Fatalf("Invalid number of groups, expected %v and received %v", 1, len(gs))
	}

	if gs[0] != "users" {
		t.Fatalf("Invalid group name: expected %v and received %v", "users", gs[0])
	}

	if err = g.Close(); err != nil {
		// @Josh this returns io.EOF
		// t.Fatal(err)
	}

	if g, err = New("./_testdata"); err != nil {
		t.Fatal(err)
	}

	if gs, err = g.Set("0", "admins"); err != nil {
		t.Fatalf("Error encountered while setting: %v", err)
	}

	if len(gs) != 2 {
		t.Fatalf("Invalid number of groups, expected %v and received %v", 2, len(gs))
	}

	if !g.Has("0", "users") {
		t.Fatal("users group does not exist when it should")
	}

	if !g.Has("0", "admins") {
		t.Fatal("admins group does not exist when it should")
	}
}
