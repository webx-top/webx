package webx

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/webx-top/echo"
)

/**
在echo框架的group.go中添加代码：

func (g *Group) URL(h Handler, params ...interface{}) string {
	return g.echo.URL(h, params...)
}

func (g *Group) SetRenderer(r Renderer) {
	g.echo.renderer = r
}

func (g *Group) Hook(h http.HandlerFunc) {
	g.echo.hook = h
}

func (g *Group) Any(path string, h Handler) {
	g.echo.Any(path, h)
}

func (g *Group) Match(methods []string, path string, h Handler) {
	g.echo.Match(methods, path, h)
}
*/
type Webxer interface {
	URL(echo.Handler, ...interface{}) string
	SetRenderer(echo.Renderer)
	Hook(http.HandlerFunc)
	Use(...echo.Middleware)
	Any(string, echo.Handler)
	Match([]string, string, echo.Handler)
	Trace(string, echo.Handler)
	WebSocket(string, echo.HandlerFunc)
	Static(string, string)
	ServeDir(string, string)
	ServeFile(string, string)
	Group(string, ...echo.Middleware) *echo.Group
}

func NewApp(name string, domain string, s *Server, middlewares ...echo.Middleware) (a *App) {
	a = &App{
		Server:      s,
		Name:        name,
		Domain:      domain,
		controllers: make(map[string]*Wrapper),
	}
	if a.Domain == "" {
		var prefix string
		if name != "" {
			prefix = `/` + name
			a.Dir = prefix + `/`
		} else {
			a.Dir = `/`
		}
		a.Url = a.Dir
		if s.Url != `/` {
			a.Url = strings.TrimSuffix(s.Url, `/`) + a.Url
		}
		a.Group = s.Echo.Group(prefix, s.DefaultMiddlewares...)
		a.Group.Use(middlewares...)
	} else {
		e := echo.New(s.InitContext)
		if s.DefaultHook != nil {
			e.Hook(s.DefaultHook)
		}
		e.Use(s.DefaultMiddlewares...)
		e.Use(middlewares...)
		if s.TemplateEngine != nil {
			e.SetRenderer(s.TemplateEngine)
		}
		a.Handler = e
		a.Url = `http://` + a.Domain + `/`
		a.Dir = `/`
	}
	return
}

type App struct {
	*Server
	*echo.Group  //没有指定域名时有效
	http.Handler //指定域名时有效
	Name         string
	Domain       string
	controllers  map[string]*Wrapper
	Url          string
	Dir          string
}

func (a *App) G() *echo.Group {
	return a.Group
}

func (a *App) E() *echo.Echo {
	return a.Handler.(*echo.Echo)
}

//注册路由：app.R(`/index`,Index.Index,"GET","POST")
func (a *App) R(path string, h HandlerFunc, methods ...string) *App {
	if len(methods) < 1 {
		methods = append(methods, "GET")
	}
	_, ctl, act := a.Server.URL.Set(path, h)
	a.Webx().Match(methods, path, func(ctx echo.Context) error {
		c := X(ctx)
		c.Init(a, nil, ctl, act)
		return h(c)
	})
	return a
}

func (a *App) Webx() Webxer {
	if a.Group != nil {
		return a.G()
	}
	return a.E()
}

//获取控制器
func (a *App) C(name string) (c interface{}) {
	c, _ = a.controllers[name]
	return
}

//登记控制器
func (a *App) RC(c interface{}) *Wrapper {
	name := fmt.Sprintf("%T", c) //example: *controller.Index
	if name[0] == '*' {
		name = name[1:]
	}
	wr := &Wrapper{
		Controller: c,
		Webx:       a.Webx(),
		App:        a,
	}
	if _, ok := c.(Initer); ok {
		_, wr.HasBefore = c.(Before)
		_, wr.HasAfter = c.(After)
	} else {
		if hf, ok := c.(BeforeHandler); ok {
			wr.BeforeHandler = hf.Before
		}
		if hf, ok := c.(AfterHandler); ok {
			wr.AfterHandler = hf.After
		}
	}
	//controller.Index
	a.controllers[name] = wr
	return wr
}
