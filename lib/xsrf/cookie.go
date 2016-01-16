package xsrf

import (
	"github.com/webx-top/echo"
	"github.com/webx-top/webx/lib/codec"
	"github.com/webx-top/webx/lib/com"
)

type CookieStorage struct {
	Prefix  string
	Secret  string
	Expires int64
	codec.Codec
}

func (c *CookieStorage) Get(key string, ctx *echo.Context) string {
	var val string
	if res, err := ctx.Request().Cookie(c.Prefix + key); err == nil && res.Value != "" {
		res.Value, _ = com.UrlDecode(res.Value)
		val = c.Codec.Decode(res.Value, c.Secret)
	}
	return val
}

func (c *CookieStorage) Set(key, val string, ctx *echo.Context) {
	val = c.Codec.Encode(val, c.Secret)
	val = com.UrlEncode(val)
	cookie := com.NewCookie(c.Prefix+key, val, c.Expires, "", "", false, true)
	ctx.Response().Header().Set("Set-Cookie", cookie.String())
}

func (c *CookieStorage) Valid(key, val string, ctx *echo.Context) bool {
	return c.Get(key, ctx) == val
}
