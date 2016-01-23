package session

import (
	"net/http"
	"strconv"
)

var DefaultName = "XSESSION"

func NewSession(engine string, setting interface{}, req *http.Request, resp http.ResponseWriter) Session {
	store := StoreEngine(engine, setting)
	return NewMySession(store, DefaultName, req, resp)
}

func NewMySession(store Store, name string, req *http.Request, resp http.ResponseWriter) Session {
	return &session{name, req, store, nil, false, resp}
}

func StoreEngine(engine string, setting interface{}) (store Store) {
	switch engine {
	case `file`:
		s := setting.(map[string]string)
		path, _ := s["path"]
		key, _ := s["key"]
		store = NewFilesystemStore(path, []byte(key))
	case `redis`:
		s := setting.(map[string]string)
		sizeStr, _ := s["size"]
		network, _ := s["network"]
		address, _ := s["address"]
		password, _ := s["password"]
		key, _ := s["key"]
		size, _ := strconv.Atoi(sizeStr)
		if size < 1 {
			size = 10
		}
		var err error
		store, err = NewRedisStore(size, network, address, password, []byte(key))
		if err != nil {
			panic(err)
		}
	case `cookie`:
		fallthrough
	default:
		s := setting.(string)
		store = NewCookieStore([]byte(s))
	}
	return
}
