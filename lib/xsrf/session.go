package xsrf

import (
	"github.com/webx-top/echo"
	X "github.com/webx-top/webx"
)

type SessionStorage struct {
}

func (c *SessionStorage) Get(key string, ctx echo.Context) string {
	s := X.X(ctx).Session()
	if s == nil {
		return ""
	}
	val, _ := s.Get(key).(string)
	return val
}

func (c *SessionStorage) Set(key, val string, ctx echo.Context) {
	s := X.X(ctx).Session()
	if s == nil {
		return
	}
	s.Set(key, val)
	s.Save()
}

func (c *SessionStorage) Valid(key, val string, ctx echo.Context) bool {
	return c.Get(key, ctx) == val
}
