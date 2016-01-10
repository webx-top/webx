package webx

import (
	"net/http"
	"strings"

	"bitbucket.org/admpub/webx/lib/pprof"
	"bitbucket.org/admpub/webx/lib/tplex"
	"github.com/admpub/echo"
	mw "github.com/admpub/echo/middleware"
	"github.com/gorilla/context"
)

func webxHeader() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Response().Header().Set(`Server`, `webx v`+VERSION)
			return h(c)
		}
	}
}

func NewServer(name string, hook http.HandlerFunc, middlewares ...echo.Middleware) (s *Server) {
	s = &Server{
		Name:               name,
		Apps:               make(map[string]*App),
		apps:               make(map[string]*App),
		DefaultMiddlewares: []echo.Middleware{webxHeader(), mw.Logger(), mw.Recover()},
		DefaultHook:        hook,
		TemplateDir:        `template`,
		Echo:               echo.New(),
	}
	s.Echo.Hook(s.DefaultHook)
	s.Echo.Use(middlewares...)
	s.Echo.Use(s.DefaultMiddlewares...)
	servs.Set(name, s)
	return
}

type Server struct {
	*echo.Echo
	Name               string
	Apps               map[string]*App //域名关联
	apps               map[string]*App //名称关联
	DefaultMiddlewares []echo.Middleware
	DefaultHook        http.HandlerFunc
	TemplateEngine     *tplex.TemplateEx
	TemplateDir        string
}

func (s *Server) SetHook(hook http.HandlerFunc) *Server {
	s.DefaultHook = hook
	s.Echo.Hook(hook)
	for _, app := range s.apps {
		app.Webx().Hook(hook)
	}
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	app, ok := s.Apps[r.Host]
	if !ok || app.Handler == nil {
		h = s.Echo
	} else {
		h = app.Handler
	}

	if h != nil {
		h.ServeHTTP(w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (s *Server) NewApp(name string, middlewares ...echo.Middleware) *App {
	r := strings.Split(name, "@") //blog@www.blog.com
	domain := ""
	if len(r) > 1 {
		name = r[0]
		domain = r[1]
	}
	a := NewApp(name, domain, s, middlewares...)
	if domain != "" {
		s.Apps[domain] = a
	}
	s.apps[name] = a
	return a
}

func (s *Server) InitTmpl(tmplDir ...string) *Server {
	if s.TemplateEngine != nil {
		s.TemplateEngine.Close()
	}
	if len(tmplDir) > 0 {
		s.TemplateEngine = tplex.New(tmplDir[0])
	} else {
		s.TemplateEngine = tplex.New(s.TemplateDir)
	}
	s.TemplateEngine.InitMgr(true, true)
	s.Echo.SetRenderer(s.TemplateEngine)
	return s
}

//启用pprof
func (s *Server) Pprof() *Server {
	pprof.Wrapper(s.Echo)
	return s
}

//开关debug模式
func (s *Server) Debug(on bool) *Server {
	s.Echo.SetDebug(on)
	return s
}

func (s *Server) Run(addr ...string) {
	http.ListenAndServe(strings.Join(addr, ":"), context.ClearHandler(s))
}

func (s *Server) App(args ...string) (a *App) {
	var name string
	if len(args) > 0 {
		name = args[0]
		if ap, ok := s.apps[name]; ok {
			a = ap
			return
		}
	}
	return s.NewApp(name)
}
