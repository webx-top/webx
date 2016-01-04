package webx

import (
	"net/http"
	"strings"

	"bitbucket.org/admpub/webx/lib/tplex"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

func NewServer(name string) (s *Server) {
	s = &Server{
		Name:               name,
		Apps:               make(map[string]*App),
		apps:               make(map[string]*App),
		DefaultMiddlewares: []echo.Middleware{mw.Logger(), mw.Recover()},
		TemplateDir:        "template",
	}
	servs.Set(name, s)
	return
}

type Server struct {
	Name               string
	Apps               map[string]*App //域名关联
	apps               map[string]*App //名称关联
	DefaultMiddlewares []echo.Middleware
	TemplateEngine     *tplex.TemplateEx
	TemplateDir        string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app, ok := s.Apps[r.Host]
	if !ok {
		app, ok = s.Apps["*"]
	}

	if ok && app.Handler != nil {
		app.Handler.ServeHTTP(w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (s *Server) New(domain string, middlewares ...echo.Middleware) *echo.Echo {
	e := echo.New()
	e.Use(middlewares...)
	e.Use(s.DefaultMiddlewares...)
	r := strings.Split(domain, "@") //blog@www.blog.com
	name := ""
	if len(r) > 1 {
		name = r[0]
		domain = r[1]
	} else {
		name = domain
	}
	a := NewApp(name, e)
	s.Apps[domain] = a
	s.apps[name] = a
	if s.TemplateEngine != nil {
		e.SetRenderer(s.TemplateEngine)
	}
	return e
}

func (s *Server) Template(tmplDir ...string) *Server {
	if s.TemplateEngine != nil {
		s.TemplateEngine.Close()
	}
	if len(tmplDir) > 0 {
		s.TemplateEngine = tplex.New(tmplDir[0])
	} else {
		s.TemplateEngine = tplex.New(s.TemplateDir)
	}
	s.TemplateEngine.InitMgr(true, true)
	return s
}

func (s *Server) Run(addr ...string) {
	http.ListenAndServe(strings.Join(addr, ":"), s)
}

func (s *Server) App(args ...string) (a *App) {
	if len(args) > 0 {
		if ap, ok := s.apps[args[0]]; ok {
			a = ap
			return
		}
	}
	return
}
