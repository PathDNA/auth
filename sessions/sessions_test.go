package sessions

import (
	"os"
	"testing"
)

const (
	testUser1 = "TEST_USER_1"
	testUser2 = "TEST_USER_2"
	testUser3 = "TEST_USER_3"
)

func TestSessions(t *testing.T) {
	var (
		s   *Sessions
		err error
	)

	s = New("./test_data")
	defer os.RemoveAll("./test_data")

	var tu1t, tu1k string
	tu1t, tu1k = s.New(testUser1)

	var tu2t, tu2k string
	tu2t, tu2k = s.New(testUser2)

	var tu3t, tu3k string
	tu3t, tu3k = s.New(testUser3)

	var mu string
	if mu, err = s.Get(tu1t, tu1k); err != nil {
		t.Fatal(err)
	} else if mu != testUser1 {
		t.Fatalf("invalid user match, expected %s and received %s", testUser1, mu)
	}

	if mu, err = s.Get(tu2t, tu2k); err != nil {
		t.Fatal(err)
	} else if mu != testUser2 {
		t.Fatalf("invalid user match, expected %s and received %s", testUser2, mu)
	}

	if mu, err = s.Get(tu3t, tu3k); err != nil {
		t.Fatal(err)
	} else if mu != testUser3 {
		t.Fatalf("invalid user match, expected %s and received %s", testUser3, mu)
	}

	if err = s.Close(); err != nil {
		t.Fatal(err)
	}

	// Re-open sessions from snapshot
	s = New("./test_data")

	// Make sure the values still match

	if mu, err = s.Get(tu1t, tu1k); err != nil {
		t.Fatal(err)
	} else if mu != testUser1 {
		t.Fatalf("invalid user match, expected %s and received %s", testUser1, mu)
	}

	if mu, err = s.Get(tu2t, tu2k); err != nil {
		t.Fatal(err)
	} else if mu != testUser2 {
		t.Fatalf("invalid user match, expected %s and received %s", testUser2, mu)
	}

	if mu, err = s.Get(tu3t, tu3k); err != nil {
		t.Fatal(err)
	} else if mu != testUser3 {
		t.Fatalf("invalid user match, expected %s and received %s", testUser3, mu)
	}
}
