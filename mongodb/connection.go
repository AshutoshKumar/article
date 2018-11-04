package mongo

import (
	"github.com/article/config"
	"gopkg.in/mgo.v2"
)

// Session ...
type Session struct {
	session *mgo.Session
}

// NewSession ...
func NewSession(config *config.MongoConfig) (*Session, error) {
	session, err := mgo.Dial(config.Ip)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	return &Session{session}, err
}

// Copy ...
func (s *Session) Copy() *mgo.Session {
	return s.session.Copy()
}

// Close ...
func (s *Session) Close() {
	if s.session != nil {
		s.session.Close()
	}
}

// DropDatabase ...
func (s *Session) DropDatabase(db string) error {
	if s.session != nil {
		return s.session.DB(db).DropDatabase()
	}
	return nil
}
