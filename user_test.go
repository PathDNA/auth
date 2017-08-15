package auth

import (
	"math/big"
	"reflect"
	"testing"
)

func TestUserAndProfile(t *testing.T) {
	type Profile struct {
		Name  string `json:"name,omitempty"`
		Phone string `json:"phone,omitempty"`
	}

	u := &User{
		ID:     big.NewInt(1),
		Status: StatusActive,

		Username: "gbusters",
		Password: "who are you gonna call",

		Profile: &Profile{
			Name:  "Ghost Busters",
			Phone: "1-800-555-2368",
		},
	}

	b, err := u.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", b)

	var p Profile

	if _, err := UnmarshalUser(b, p); err != ErrProfileNotPtr {
		t.Fatalf("expected ErrProfileNotPtr, got %v", err)
	}

	nu, err := UnmarshalUser(b, &p)
	if err != nil {
		t.Fatal(err)
	}

	if nu.CreatedTS == 0 {
		t.Fatal("created ts wasn't set")
	} else {
		nu.CreatedTS = 0 // need to be 0 for DeepEqual
	}

	if !reflect.DeepEqual(u, nu) {
		t.Logf("u != nu\n%#+v\n%#+v", u, nu)
	}
}
