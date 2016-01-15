package xsrf

import (
	"github.com/gorilla/context"
	"github.com/webx-top/webx/echo"
	"github.com/webx-top/webx/lib/middleware/session"
)

type SessionStorage struct {
	StoreEngine string
	Setting     interface{}
}

func (c *SessionStorage) Init(ctx *echo.Context) Session {
	s := session.NewSession(c.StoreEngine, c.Setting, ctx.Request(), ctx.Response().Writer())
	context.Set(r, `XsrfSession`, s)
	return s
}

func (c *SessionStorage) Get(key string, ctx *echo.Context) string {
	s, ok := context.Get(r, `XsrfSession`).(Session)
	if !ok {
		s = c.Init(ctx)
	}
	return s.Get(key).(string)
}

func (c *SessionStorage) Set(key, val string, ctx *echo.Context) {
	s, ok := context.Get(r, `XsrfSession`).(Session)
	if !ok {
		s = c.Init(ctx)
	}
	s.Set(key, val)
}

func (c *SessionStorage) Valid(key, val string, ctx *echo.Context) bool {
	return c.Get(key, ctx) == val
}
