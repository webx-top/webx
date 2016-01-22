package webx

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	"github.com/webx-top/echo"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/webx/lib/pprof"
	"github.com/webx-top/webx/lib/tplex"
	"github.com/webx-top/webx/lib/tplfunc"
)

func webxHeader() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
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
		Url:                `/`,
		InitializeContext: func(resp *echo.Response, e *echo.Echo) interface{} {
			return NewContext(echo.NewContext(nil, resp, e))
		},
	}
	s.Echo = echo.New(s.InitializeContext)
	s.URL = NewURL(name, s)
	s.Echo.Hook(s.DefaultHook)
	s.Echo.Use(s.DefaultMiddlewares...)
	s.Echo.Use(middlewares...)
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
	Url                string
	*URL
	InitializeContext func(*echo.Response, *echo.Echo) interface{}
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
	s.Echo.Logger().Info(`Server "%v" has been launched.`, s.Name)
	err := http.ListenAndServe(strings.Join(addr, ":"), context.ClearHandler(s))
	if err != nil {
		s.Echo.Logger().Error(err)
	}
	s.Echo.Logger().Info(`Server "%v" has been closed.`, s.Name)
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

func (s *Server) FuncMap() (f template.FuncMap) {
	f = tplfunc.TplFuncMap
	f["UrlFor"] = s.URL.BuildByPath
	f["Url"] = s.URL.Build
	f["RootUrl"] = func(p ...string) string {
		if len(p) > 0 {
			return s.Url + p[0]
		}
		return s.Url
	}
	return
}

func (s *Server) Static(absPath string, urlPath string, f ...*template.FuncMap) *tplfunc.Static {
	st := tplfunc.NewStatic(absPath, urlPath)
	if len(f) > 0 {
		*f[0] = st.Register(*f[0])
	}
	return st
}
