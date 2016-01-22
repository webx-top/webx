package webx

import (
	"bytes"
	"io/ioutil"
	"strings"
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/webx/lib/com"
	"github.com/webx-top/webx/lib/cookie"
	sessionMW "github.com/webx-top/webx/lib/middleware/session"
)

func NewContext(s *Server, c echo.Context) *Context {
	return &Context{
		Context: c,
		Server:  s,
	}
}

type Context struct {
	*Server
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

func (c *Context) Cookie(key string, value string) *cookie.Cookie {
	liftTime := c.Server.CookieExpires
	sPath := "/"
	domain := c.Server.CookieDomain
	secure := c.Server.CookieSecure
	httpOnly := c.Server.CookieHttpOnly
	return cookie.New(c.Server.CookiePrefix+key, value, liftTime, sPath, domain, secure, httpOnly)
}

func (c *Context) GetCookie(key string) string {
	var val string
	if res, err := c.Request().Cookie(c.Server.CookiePrefix + key); err == nil && res.Value != "" {
		val, _ = com.UrlDecode(res.Value)
	}
	return val
}

func (c *Context) SetCookie(key, val string, args ...interface{}) {
	val = com.UrlEncode(val)
	cookie := c.Cookie(key, val)
	switch len(args) {
	case 5:
		httpOnly, _ := args[4].(bool)
		cookie.HttpOnly(httpOnly)
		fallthrough
	case 4:
		secure, _ := args[3].(bool)
		cookie.Secure(secure)
		fallthrough
	case 3:
		domain, _ := args[2].(string)
		cookie.Domain(domain)
		fallthrough
	case 2:
		path, _ := args[1].(string)
		cookie.Path(path)
		fallthrough
	case 1:
		var liftTime int64
		switch args[0].(type) {
		case int:
			liftTime = int64(args[0].(int))
		case int64:
			liftTime = args[0].(int64)
		case time.Duration:
			liftTime = int64(args[0].(time.Duration))
		}
		cookie.Expires(liftTime)
	}
	cookie.Send(c)
}

func (c *Context) SetSecCookie(key string, value interface{}) {
	if c.Server.Codec == nil {
		val, _ := value.(string)
		c.SetCookie(key, val)
		return
	}
	encoded, err := c.Server.Codec.Encode(key, value)
	if err != nil {
		c.X().Echo().Logger().Error(err)
	} else {
		c.SetCookie(key, encoded)
	}
}

func (c *Context) GetSecCookie(key string) (value interface{}) {
	cookieValue := c.GetCookie(key)
	if cookieValue != "" && c.Server.Codec != nil {
		err := c.Server.Codec.Decode(key, cookieValue, &value)
		if err != nil {
			c.X().Echo().Logger().Error(err)
		}
	}
	return
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
