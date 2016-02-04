package session

import (
	"net/http"

	ss "github.com/webx-top/webx/lib/session/engine/gorilla"
	in "github.com/webx-top/webx/lib/session/ssi"
)

func NewSession(engine string, setting interface{}, req *http.Request, resp http.ResponseWriter) in.Session {
	return ss.NewSession(engine, setting, req, resp)
}

type Store interface {
	ss.Store
}

func NewMySession(store ss.Store, name string, req *http.Request, resp http.ResponseWriter) in.Session {
	return ss.NewMySession(store, name, req, resp)
}

func StoreEngine(engine string, setting interface{}) (store Store) {
	return ss.StoreEngine(engine, setting)
}
