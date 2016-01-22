package webx

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/webx/lib/com"
	"github.com/webx-top/webx/lib/cookie"
	sessionMW "github.com/webx-top/webx/lib/middleware/session"
)

func NewContext(c echo.Context) *Context {
	return &Context{
		Context: c,
	}
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

func (c *Context) Cookie(key string, value string, args ...interface{}) *cookie.Cookie {
	return cookie.New(key, value, args...)
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
	c.Cookie(key, val, args...).Send(c)
}

func (c *Context) Body() ([]byte, error) {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil, err
	}

	c.Request().Body.Close()
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body))

	return body, nil
}

func (c *Context) IP() string {
	proxy := []string{}
	if ips := c.Request().Header.Get("X-Forwarded-For"); ips != "" {
		proxy = strings.Split(ips, ",")
	}
	if len(proxy) > 0 && proxy[0] != "" {
		return proxy[0]
	}
	ip := strings.Split(c.Request().RemoteAddr, ":")
	if len(ip) > 0 {
		if ip[0] != "[" {
			return ip[0]
		}
	}
	return "127.0.0.1"
}

func (c *Context) IsAjax() bool {
	return c.Request().Header.Get("X-Requested-With") == "XMLHttpRequest"
}
