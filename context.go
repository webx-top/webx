package webx

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/webx/lib/com"
	sessionMW "github.com/webx-top/webx/lib/middleware/session"
)

func NewContext() *Context {
	return &Context{}
}

type Context struct {
	echo.Context
	session sessionMW.Session
}

func (c *Context) Session() sessionMW.Session {
	if c.session == nil {
		c.session = sessionMW.Default(c)
	}
	return c.session
}

func (c *Context) SetSession(key string, val interface{}) {
	s := c.Session()
	s.Set(key, val)
	s.Save()
}

func (c *Context) GetSession(key string) interface{} {
	return c.Session().Get(key)
}

func (c *Context) GetCookie(key string) string {
	var val string
	if res, err := c.Request().Cookie(key); err == nil && res.Value != "" {
		val, _ = com.UrlDecode(res.Value)
	}
	return val
}

func (c *Context) SetCookie(key, val string, args ...interface{}) {
	val = com.UrlEncode(val)
	cookie := com.NewCookie(key, val, args...)
	c.Response().Header().Set("Set-Cookie", cookie.String())
}
