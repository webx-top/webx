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
	return echo.MiddlewareFunc(func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(ctx echo.Context) error {
			c := X.X(ctx)
			s := ss.NewMySession(store, name, ctx)
			if se, ok := interface{}(c).(Sessionser); ok {
				se.InitSession(s)
			}
			err := h.Handle(c)
			s.Save()
			return err
		})
	})
}

func Middleware(options *ssi.Options, setting interface{}) echo.MiddlewareFunc {
	store := ss.StoreEngine(options, setting)
	return Sessions(ssi.DefaultName, store)
}
