// Session implements middleware for easily using github.com/gorilla/sessions
// within echo. This package was originally inspired from the
// https://github.com/ipfans/echo-session package, and modified to provide more
// functionality
package session

import (
	"github.com/webx-top/echo"
	X "github.com/webx-top/webx"
	sessLib "github.com/webx-top/webx/lib/session"
)

func Sessions(name string, store sessLib.Store) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.IsFileServer() {
				return h(c)
			}
			s := sessLib.NewMySession(store, name, c.Request(), c.Response().Writer())
			X.X(c).InitSession(s)
			return h(c)
		}
	}
}

func Middleware(engine string, setting interface{}) echo.MiddlewareFunc {
	store := sessLib.StoreEngine(engine, setting)
	return Sessions(sessLib.DefaultName, store)
}
