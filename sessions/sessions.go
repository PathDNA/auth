package sessions

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/Path94/atoms"
	"github.com/missionMeteora/toolkit/errors"
	"github.com/missionMeteora/uuid"
)

const (
	// ErrSessionDoesNotExist is returned when an invalid token/key pair is presented
	ErrSessionDoesNotExist = errors.Error("session with that token/key pair does not exist")
)

const (
	// SessionTimeout (in seconds) is the ttl for sessions, an action will refresh the duration
	SessionTimeout = 60 * 60 * 12 // 12 hours
)

const (
	snapshotName = "sessions.db"
)

// New will return a new instance of sessions
func New(dir string) *Sessions {
	var s Sessions
	s.dir = dir
	s.g = uuid.NewGen()
	s.m = make(map[string]*session)
	// Load from snapshot
	s.load()
	// Start purge loop
	go s.loop()
	return &s
}

// Sessions manages sessions
type Sessions struct {
	mux atoms.RWMux

	dir string

	g *uuid.Gen
	m map[string]*session

	closed atoms.Bool
}

func (s *Sessions) loop() {
	for !s.closed.Get() {
		oldest := time.Now().Add(time.Second * -SessionTimeout).Unix()
		s.Purge(oldest)
		time.Sleep(time.Minute)
	}
}

// Purge will purge all entries oldest than the oldest value
func (s *Sessions) Purge(oldest int64) {
	s.mux.Update(func() {
		for key, ss := range s.m {
			if ss.LastAction.Load() < oldest {
				delete(s.m, key)
			}
		}
	})
}

// New will creata  new token/key pair
func (s *Sessions) New(uuid string) (token, key string) {
	var ss session
	ss.UUID = uuid
	ss.setAction()

	// Set token
	token = s.g.New().String()
	// Set key
	key = s.g.New().String()

	s.mux.Update(func() {
		s.m[getMapKey(token, key)] = &ss
	})

	return
}

// Get will retrieve the UUID associated with a provided token/key pair
func (s *Sessions) Get(token, key string) (uuid string, err error) {
	var (
		ss *session
		ok bool
	)

	s.mux.Read(func() {

		if ss, ok = s.m[getMapKey(token, key)]; !ok {
			err = ErrSessionDoesNotExist
			return
		}

		// Set uuid as session uuid
		uuid = ss.UUID
		// Set last action for session
		ss.setAction()
	})

	return
}

func (s *Sessions) load() (err error) {
	var f *os.File
	if f, err = os.Open(filepath.Join(s.dir, snapshotName)); err != nil {
		return
	}
	defer f.Close()

	return json.NewDecoder(f).Decode(&s.m)
}

func (s *Sessions) snapshot() (err error) {
	if err = os.MkdirAll(s.dir, 0744); err != nil {
		return
	}

	var f *os.File
	if f, err = os.Create(filepath.Join(s.dir, snapshotName)); err != nil {
		return
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(s.m)
}

// Close will close an instance of Sessions
func (s *Sessions) Close() (err error) {
	if !s.closed.Set(true) {
		return errors.ErrIsClosed
	}

	return s.snapshot()
}

func getMapKey(token, key string) (mapkey string) {
	return token + "::" + key
}
