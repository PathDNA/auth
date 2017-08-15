package auth

import (
	"encoding/json"
	"reflect"
	"testing"
)

type Profile struct {
	Name  string `json:"name,omitempty"`
	Phone string `json:"phone,omitempty"`
}

func TestUserAndProfile(t *testing.T) {
	u := &User{
		ID:     "1",
		Status: StatusActive,

		Username: "gbusters",
		Password: "who are you gonna call",

		Profile: &Profile{
			Name:  "Ghost Busters",
			Phone: "1-800-555-2368",
		},
	}

	b, err := json.Marshal(u)
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

	if !reflect.DeepEqual(u, nu) {
		t.Errorf("u != nu\n%#+v\n%#+v", u, nu)
	}
}
