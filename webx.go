package webx

import (
	"net/http"
	"strings"

	"bitbucket.org/admpub/webx/lib/pprof"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

var Serv *Server=NewServer("webx")
var servs *Servers=new(Servers)


//========================================Servers
type Servers map[string]*Server

func (s *Servers) Init(name string) (sv *Server) {
	s=new(map[string]*Server)
}

func (s *Servers) Get(name string) (sv *Server) {
	sv,_=(*s)[name]
	return
}

func (s *Servers) Set(name string,sv *Server) {
	(*s)[name]=sv
}


//========================================Server
func NewServer(name string) (s *Server) {
	s=&Server{
		Name:name,
		Apps:make(map[string]*App),
		apps:make(map[string]*App),
		DefaultMiddlewares:[]echo.MiddlewareFunc{mw.Logger(),mw.Recover()},
	}
	servs.Set(name, s)
	return
}

type Server struct {
	Name string
	Apps map[string]*App //域名关联
	apps map[string]*App //名称关联
	DefaultMiddlewares []echo.MiddlewareFunc
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app, ok := s.Apps[r.Host]
	if !ok {
		app, ok = s.Apps["*"]
	} 

	if ok && app.Handler != nil {
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
	a:=NewApp(name,e)
	s.Apps[domain] = a
	s.apps[name] = a
	return e
}

func (s *Server) Run(addr ...string){
	http.ListenAndServe(strings.Join(addr, ":"), s)
}


//========================================App
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
