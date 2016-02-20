package session

import (
	"github.com/webx-top/echo"
	ss "github.com/webx-top/webx/lib/session/engine/gorilla"
	in "github.com/webx-top/webx/lib/session/ssi"
)

func NewSession(options *in.Options, setting interface{}, ctx echo.Context) in.Session {
	return ss.NewSession(options, setting, ctx)
}

type Store interface {
	ss.Store
}

func NewMySession(store ss.Store, name string, ctx echo.Context) in.Session {
	return ss.NewMySession(store, name, ctx)
}

func StoreEngine(options *in.Options, setting interface{}) (store Store) {
	return ss.StoreEngine(options, setting)
}
