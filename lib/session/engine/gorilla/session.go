// Session implements middleware for easily using github.com/gorilla/sessions
// within echo. This package was originally inspired from the
// https://github.com/ipfans/echo-session package, and modified to provide more
// functionality
package session

import (
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	I "github.com/webx-top/webx/lib/session/ssi"
)

const (
	errorFormat = "[sessions] ERROR! %s\n"
)

type Store interface {
	sessions.Store
	Options(I.Options)
}

type Session struct {
	name    string
	request *http.Request
	store   Store
	session *sessions.Session
	written bool
	writer  http.ResponseWriter
}

func (s *Session) Get(key string) interface{} {
	return s.Session().Values[key]
}

func (s *Session) Set(key string, val interface{}) I.Session {
	s.Session().Values[key] = val
	s.written = true
	return s
}

func (s *Session) Delete(key string) I.Session {
	delete(s.Session().Values, key)
	s.written = true
	return s
}

func (s *Session) Clear() I.Session {
	for key := range s.Session().Values {
		if k, ok := key.(string); ok {
			s.Delete(k)
		}
	}
	return s
}

func (s *Session) AddFlash(value interface{}, vars ...string) I.Session {
	s.Session().AddFlash(value, vars...)
	s.written = true
	return s
}

func (s *Session) Flashes(vars ...string) []interface{} {
	s.written = true
	return s.Session().Flashes(vars...)
}

func (s *Session) Options(options I.Options) I.Session {
	s.Session().Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
	return s
}

func (s *Session) SetID(id string) I.Session {
	s.Session().ID = id
	return s
}

func (s *Session) Save() error {
	if s.Written() {
		e := s.Session().Save(s.request, s.writer)
		if e == nil {
			s.written = false
		}
		return e
	}
	return nil
}

func (s *Session) Session() *sessions.Session {
	if s.session == nil {
		var err error
		s.session, err = s.store.Get(s.request, s.name)
		if err != nil {
			log.Printf(errorFormat, err)
		}
	}
	return s.session
}

func (s *Session) Written() bool {
	return s.written
}
