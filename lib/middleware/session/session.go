// Session implements middleware for easily using github.com/gorilla/sessions
// within echo. This package was originally inspired from the
// https://github.com/ipfans/echo-session package, and modified to provide more
// functionality
package session

import (
	"github.com/webx-top/echo"
	X "github.com/webx-top/webx"
	ss "github.com/webx-top/webx/lib/session"
	"github.com/webx-top/webx/lib/session/ssi"
)

type Sessionser interface {
	InitSession(ssi.Session)
}

func Sessions(name string, store ss.Store) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if ctx.IsFileServer() {
				return h(ctx)
			}
			c := X.X(ctx)
			s := ss.NewMySession(store, name, c.Request(), c.Response().Writer())
			if se, ok := interface{}(c).(Sessionser); ok {
				se.InitSession(s)
			}
			err := h(c)
			s.Save()
			return err
		}
	}
}

func Middleware(options *ssi.Options, setting interface{}) echo.MiddlewareFunc {
	store := ss.StoreEngine(options, setting)
	return Sessions(ssi.DefaultName, store)
}
