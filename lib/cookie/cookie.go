package cookie

import (
	"net/http"
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/webx/lib/com"
)

func New(name string, value string, args ...interface{}) *Cookie {
	return &Cookie{
		cookie: com.NewCookie(name, value, args...),
	}
}

type Cookie struct {
	cookie *http.Cookie
	/*
		Name:     name,
		Value:    value,
		Path:     path,
		Domain:   domain,
		MaxAge:   0,
		Secure:   secure,
		HttpOnly: httpOnly,
	*/
}

func (c *Cookie) Path(p string) *Cookie {
	c.cookie.Path = p
	return c
}

func (c *Cookie) Domain(p string) *Cookie {
	c.cookie.Domain = p
	return c
}

func (c *Cookie) MaxAge(p int) *Cookie {
	c.cookie.MaxAge = p
	return c
}

func (c *Cookie) Expires(p int64) *Cookie {
	if p > 0 {
		c.cookie.Expires = time.Unix(time.Now().Unix()+p, 0)
	} else if p < 0 {
		c.cookie.Expires = time.Unix(1, 0)
	}
	return c
}

func (c *Cookie) Secure(p bool) *Cookie {
	c.cookie.Secure = p
	return c
}

func (c *Cookie) HttpOnly(p bool) *Cookie {
	c.cookie.HttpOnly = p
	return c
}

func (c *Cookie) Send(ctx echo.Context) {
	ctx.Response().Header().Set("Set-Cookie", c.cookie.String())
}
