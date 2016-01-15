package session

import (
	"strconv"

	"github.com/webx-top/echo"
)

func NewSession(engine string, setting interface{}, req *http.Request, resp http.ResponseWriter) Session {
	store := StoreEngine(engine, setting)
	return &session{"XSESSION", req, store, nil, false, resp}
}

func Middleware(engine string, setting interface{}) echo.MiddlewareFunc {
	store := StoreEngine(engine, setting)
	return Sessions("XSESSION", store)
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
