/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package webx

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/context"
	codec "github.com/gorilla/securecookie"
	"github.com/webx-top/echo"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/webx/lib/events"
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
		MaxUploadSize:      10 * 1024 * 1024,
		CookiePrefix:       "webx_" + name + "_",
		CookieHttpOnly:     true,
	}
	s.InitContext = func(resp *echo.Response, e *echo.Echo) interface{} {
		return NewContext(s, echo.NewContext(nil, resp, e))
	}

	s.CookieAuthKey = string(codec.GenerateRandomKey(32))
	s.CookieBlockKey = string(codec.GenerateRandomKey(32))
	s.SessionStoreEngine = `cookie`
	s.SessionStoreConfig = s.CookieAuthKey
	s.Codec = codec.New([]byte(s.CookieAuthKey), []byte(s.CookieBlockKey))
	s.Core = echo.New(s.InitContext)
	s.URL = NewURL(name, s)
	s.Core.Hook(s.DefaultHook)
	s.Core.Use(s.DefaultMiddlewares...)
	s.Core.Use(middlewares...)
	servs.Set(name, s)
	return
}

type Server struct {
	Core               *echo.Echo
	Name               string
	Apps               map[string]*App //域名关联
	apps               map[string]*App //名称关联
	DefaultMiddlewares []echo.Middleware
	DefaultHook        http.HandlerFunc
	TemplateEngine     *tplex.TemplateEx
	TemplateDir        string
	MaxUploadSize      int64
	CookiePrefix       string
	CookieHttpOnly     bool
	CookieAuthKey      string
	CookieBlockKey     string
	CookieExpires      int64
	CookieDomain       string
	SessionStoreEngine string
	SessionStoreConfig interface{}
	codec.Codec
	Url string
	*URL
	InitContext func(*echo.Response, *echo.Echo) interface{}
}

//初始化 加密/解密 接口
func (s *Server) InitCodec(hashKey []byte, blockKey []byte) {
	s.Codec = codec.New(hashKey, blockKey)
}

//设置钩子函数
func (s *Server) SetHook(hook http.HandlerFunc) *Server {
	s.DefaultHook = hook
	s.Core.Hook(hook)
	for _, app := range s.apps {
		app.Webx().Hook(hook)
	}
	return s
}

//HTTP服务执行入口
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var h http.Handler
	app, ok := s.Apps[r.Host]
	if !ok || app.Handler == nil {
		h = s.Core
	} else {
		h = app.Handler
	}

	if h != nil {
		h.ServeHTTP(w, r)
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

//创建新app
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

//初始化模板引擎
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
	s.Core.SetRenderer(s.TemplateEngine)
	return s
}

//启用pprof
func (s *Server) Pprof() *Server {
	pprof.Wrapper(s.Core)
	return s
}

//开关debug模式
func (s *Server) Debug(on bool) *Server {
	s.Core.SetDebug(on)
	return s
}

//运行服务
func (s *Server) Run(addr ...string) {
	defer func() {
		events.GoEvent(`webx.serverExit`, nil, func(_ bool) {})
	}()
	s.Core.Logger().Info(`Server "%v" has been launched.`, s.Name)
	err := http.ListenAndServe(strings.Join(addr, ":"), context.ClearHandler(s))
	if err != nil {
		s.Core.Logger().Error(err)
	}
	s.Core.Logger().Info(`Server "%v" has been closed.`, s.Name)
}

//已创建app实例
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

//可用全局模板函数
func (s *Server) FuncMap() (f template.FuncMap) {
	f = tplfunc.TplFuncMap
	f["AppUrlFor"] = s.URL.BuildByPath
	f["AppUrl"] = s.URL.Build
	f["RootUrl"] = func(p ...string) string {
		if len(p) > 0 {
			return s.Url + p[0]
		}
		return s.Url
	}
	return
}

//静态资源文件管理器
func (s *Server) Static(absPath string, urlPath string, f ...*template.FuncMap) *tplfunc.Static {
	st := tplfunc.NewStatic(absPath, urlPath)
	if len(f) > 0 {
		*f[0] = st.Register(*f[0])
	}
	return st
}
