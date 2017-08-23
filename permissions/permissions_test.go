package permissions

import (
	"os"
	"testing"

	"github.com/missionMeteora/toolkit/errors"
)

const (
	testErrCannot = errors.Error("group not allowed to perform action they should be able to")
	testErrCan    = errors.Error("group allowed to perform action they should not be able to")
)

func TestPermissions(t *testing.T) {
	var (
		p   *Permissions
		err error
	)

	if p, err = New("./_testdata"); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll("./_testdata")

	if err = p.Set("0", "users", PermissionRead); err != nil {
		t.Fatal(err)
	}

	if err = p.Set("0", "admins", PermissionReadWrite); err != nil {
		t.Fatal(err)
	}

	if err = p.Set("0", "writers", PermissionWrite); err != nil {
		t.Fatal(err)
	}

	testPerms(p, t)

	if err = p.Close(); err != nil {
		t.Error(err)
	}

	if p, err = New("./_testdata"); err != nil {
		t.Fatal(err)
	}

	testPerms(p, t)
}

func testPerms(p *Permissions, t *testing.T) {
	if p.Can("0", ActionWrite, "users") {
		t.Fatal(testErrCan)
	}

	if !p.Can("0", ActionWrite, "admins") {
		t.Fatal(testErrCannot)
	}

	if !p.Can("0", ActionWrite, "writers") {
		t.Fatal(testErrCannot)
	}

	if !p.Can("0", ActionWrite, "users", "admins") {
		t.Fatal(testErrCannot)
	}

	if !p.Can("0", ActionRead, "users") {
		t.Fatal(testErrCannot)
	}

	if !p.Can("0", ActionRead, "admins") {
		t.Fatal(testErrCannot)
	}

	if p.Can("0", ActionRead, "writers") {
		t.Fatal(testErrCan)
	}
}
