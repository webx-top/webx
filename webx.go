package webx

import (
	"net/http"
	"strings"

	"bitbucket.org/admpub/webx/lib/pprof"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

func NewServer(name string) (s *Server) {
	s=&Server{
		Name:name,
		Apps:make(map[string]*App),
		apps:make(map[string]*App),
		DefaultMiddlewares:[]echo.MiddlewareFunc{mw.Logger(),mw.Recover()},
	}
	return 
}

type Server struct {
	Name string
	Apps map[string]*App
	apps map[string]*App
	DefaultMiddlewares []echo.MiddlewareFunc
}

func NewApp(name string,e *echo.Echo) (a *App) {
	a=&App{
		Name:name,
		Handler:e,
	}
	return
}

type App struct {
	Name string
	http.Handler
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if app, ok := s.Apps[r.Host]; ok && app.Handler != nil {
		app.Handler.ServeHTTP(w, r)
	} else if app, ok := s.Apps["*"]; ok && app.Handler != nil {
		app.Handler.ServeHTTP(w, r) {
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (s *Server) New(domain string,middlewares ...echo.MiddlewareFunc) *echo.Echo {
	e := echo.New()
	e.Use(middlewares...)
	e.Use(s.DefaultMiddlewares...)
	r:=strings.Split(domain, "@") //blog@www.blog.com
	name:=""
	if len(r)>1 {
		name=r[0]
		domain=r[1]
	}else{
		name=domain
	}
	a:=NewApp(domain,e)
	s.Apps[domain] = a
	s.apps[name] = a
	return e
}

func (s *Server) Run(addr ...string){
	http.ListenAndServe(strings.Join(addr, ":"), s)
}

func main() {
	server:=NewServer("webx")
	server.Run(":8080")
}
