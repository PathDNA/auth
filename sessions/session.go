package sessions

import (
	"time"

	"github.com/PathDNA/atoms"
)

type session struct {
	UUID string `json:"uuid"`
	// Last action taken for this session
	LastAction atoms.Int64 `json:"lastAction"`
}

func (s *session) setAction() {
	s.LastAction.Store(time.Now().Unix())
}
