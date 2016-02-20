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
	"strings"

	codec "github.com/gorilla/securecookie"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/webx/lib/events"
	"github.com/webx-top/webx/lib/pprof"
	"github.com/webx-top/webx/lib/tplex"
	"github.com/webx-top/webx/lib/tplfunc"
)

func webxHeader() echo.MiddlewareFunc {
	return echo.MiddlewareFunc(func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.Response().Header().Set(`Server`, `webx v`+VERSION)
			return h.Handle(c)
		})
	})
}

func NewServer(name string, middlewares ...echo.Middleware) (s *Server) {
	s = &Server{
		Name:               name,
		Apps:               make(map[string]*App),
		apps:               make(map[string]*App),
		DefaultMiddlewares: []echo.Middleware{webxHeader(), mw.Log(), mw.Recover()},
		TemplateDir:        `template`,
		Url:                `/`,
		MaxUploadSize:      10 * 1024 * 1024,
		CookiePrefix:       "webx_" + name + "_",
		CookieHttpOnly:     true,
	}
	s.InitContext = func(e *echo.Echo) interface{} {
		return NewContext(s, echo.NewContext(nil, nil, e))
	}

	s.CookieAuthKey = string(codec.GenerateRandomKey(32))
	s.CookieBlockKey = string(codec.GenerateRandomKey(32))
	s.SessionStoreEngine = `cookie`
	s.SessionStoreConfig = s.CookieAuthKey
	s.Codec = codec.New([]byte(s.CookieAuthKey), []byte(s.CookieBlockKey))
	s.Core = echo.NewWithContext(s.InitContext)
	s.URL = NewURL(name, s)
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
	TemplateEngine     tplex.TemplateEx
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
	InitContext func(*echo.Echo) interface{}
}

// 初始化 加密/解密 接口
func (s *Server) InitCodec(hashKey []byte, blockKey []byte) {
	s.Codec = codec.New(hashKey, blockKey)
}

// HTTP服务执行入口
func (s *Server) ServeHTTP(r engine.Request, w engine.Response) {
	var h *echo.Echo
	app, ok := s.Apps[r.Host()]
	if !ok || app.Handler == nil {
		h = s.Core
	} else {
		h = app.Handler
	}

	if h != nil {
		h.ServeHTTP(r, w)
	} else {
		w.NotFound()
	}
}

// 创建新app
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

// 重置模板引擎
func (s *Server) ResetTmpl(args ...interface{}) *Server {
	if s.TemplateEngine != nil {
		s.TemplateEngine.Close()
	}
	s.TemplateEngine = s.InitTmpl(args...)
	s.Core.SetRenderer(s.TemplateEngine)
	return s
}

// 初始化模板引擎
func (s *Server) InitTmpl(args ...interface{}) (tmplEng tplex.TemplateEx) {
	var tmplDir, engine string
	var cachedContent, reloadTmpl = true, true
	switch len(args) {
	case 4:
		reloadTmpl, _ = args[3].(bool)
		fallthrough
	case 3:
		cachedContent, _ = args[2].(bool)
		fallthrough
	case 2:
		engine, _ = args[1].(string)
		fallthrough
	case 1:
		tmplDir, _ = args[0].(string)
	}
	if tmplDir == `` {
		tmplDir = s.TemplateDir
	}
	tmplEng = tplex.Create(engine, tmplDir)
	tmplEng.Init(cachedContent, reloadTmpl)
	return
}

// 启用pprof
func (s *Server) Pprof() *Server {
	pprof.Wrapper(s.Core)
	return s
}

// 开关debug模式
func (s *Server) Debug(on bool) *Server {
	s.Core.SetDebug(on)
	return s
}

// 运行服务
func (s *Server) Run(args ...interface{}) {
	var eng engine.Engine
	var arg interface{}
	if len(args) > 0 {
		arg = args[0]
	}
	switch arg.(type) {
	case string:
		eng = standard.New(arg.(string))
	case engine.Engine:
		eng = args[0].(engine.Engine)
	default:
		eng = standard.New(`:80`)
	}
	defer func() {
		events.GoEvent(`webx.serverExit`, nil, func(_ bool) {})
	}()
	s.Core.Logger().Infof(`Server "%v" has been launched.`, s.Name)

	eng.SetHandler(s.ServeHTTP)
	eng.SetLogger(s.Core.Logger())
	eng.Start()

	s.Core.Logger().Infof(`Server "%v" has been closed.`, s.Name)
}

// 已创建app实例
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

// 可用全局模板函数
func (s *Server) FuncMap() (f map[string]interface{}) {
	f = map[string]interface{}{}
	for k, v := range tplfunc.TplFuncMap {
		f[k] = v
	}
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

// 静态资源文件管理器
func (s *Server) Static(absPath string, urlPath string, f ...*map[string]interface{}) *tplfunc.Static {
	st := tplfunc.NewStatic(absPath, urlPath)
	if len(f) > 0 {
		*f[0] = st.Register(*f[0])
	}
	return st
}
