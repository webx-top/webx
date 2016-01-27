package xsrf

import (
	"github.com/webx-top/echo"
	X "github.com/webx-top/webx"
)

type CookieStorage struct {
}

func (c *CookieStorage) Get(key string, ctx echo.Context) string {
	return X.X(ctx).GetCookie(key)
}

func (c *CookieStorage) Set(key, val string, ctx echo.Context) {
	X.X(ctx).SetCookie(key, val)
}

func (c *CookieStorage) Valid(key, val string, ctx echo.Context) bool {
	return c.Get(key, ctx) == val
}
