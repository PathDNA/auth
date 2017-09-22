package sessions

import (
	"time"

	"github.com/Path94/atoms"
)

type session struct {
	uuid       string
	lastAction atoms.Int64
}

func (s *session) setAction() {
	s.lastAction.Store(time.Now().Unix())
}
