package xsrf

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/webx/lib/middleware/session"
)

//依赖于session.Middleware(engine string, setting interface{})中间件
type SessionStorage struct {
}

func (c *SessionStorage) Get(key string, ctx echo.Context) string {
	s := session.Default(ctx)
	if s == nil {
		return ""
	}
	return s.Get(key).(string)
}

func (c *SessionStorage) Set(key, val string, ctx echo.Context) {
	s := session.Default(ctx)
	if s == nil {
		return
	}
	s.Set(key, val)
	s.Save()
}

func (c *SessionStorage) Valid(key, val string, ctx echo.Context) bool {
	return c.Get(key, ctx) == val
}
