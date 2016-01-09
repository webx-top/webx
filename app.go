package webx

import (
	"fmt"
	"net/http"

	"github.com/admpub/echo"
)

/**
在echo框架的group.go中添加代码：

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
		controllers: make(map[string]*Controller),
	}
	if a.Domain == "" {
		var prefix string
		if name != "" {
			prefix = `/` + name
		}
		a.Group = s.Echo.Group(prefix, middlewares...)
		a.Group.Use(s.DefaultMiddlewares...)
		if s.DefaultHook != nil {
			a.Group.Hook(s.DefaultHook)
		}
	} else {
		e := echo.New()
		if s.DefaultHook != nil {
			e.Hook(s.DefaultHook)
		}
		e.Use(middlewares...)
		e.Use(s.DefaultMiddlewares...)
		if s.TemplateEngine != nil {
			e.SetRenderer(s.TemplateEngine)
		}
		a.Handler = e
	}
	return
}

type Controller struct {
	Before     echo.HandlerFunc
	After      echo.HandlerFunc
	Controller interface{}
	Webx       Webxer
}

//注册路由：Controller.R(`/index`,Index.Index,"GET","POST")
func (a *Controller) R(path string, h echo.HandlerFunc, methods ...string) *Controller {
	if len(methods) < 1 {
		methods = append(methods, "GET")
	}
	if a.Before != nil && a.After != nil {
		a.Webx.Match(methods, path, func(c *echo.Context) error {
			c.Set(`Exit`, false)
			if err := a.Before(c); err != nil {
				return err
			}
			if exit, _ := c.Get(`Exit`).(bool); exit {
				return nil
			}
			if err := h(c); err != nil {
				return err
			}
			if exit, _ := c.Get(`Exit`).(bool); exit {
				return nil
			}
			return a.After(c)
		})
	} else if a.Before != nil {
		a.Webx.Match(methods, path, func(c *echo.Context) error {
			c.Set(`Exit`, false)
			if err := a.Before(c); err != nil {
				return err
			}
			if exit, _ := c.Get(`Exit`).(bool); exit {
				return nil
			}
			return h(c)
		})
	} else if a.After != nil {
		a.Webx.Match(methods, path, func(c *echo.Context) error {
			c.Set(`Exit`, false)
			if err := h(c); err != nil {
				return err
			}
			if exit, _ := c.Get(`Exit`).(bool); exit {
				return nil
			}
			return a.After(c)
		})
	} else {
		a.Webx.Match(methods, path, func(c *echo.Context) error {
			c.Set(`Exit`, false)
			return h(c)
		})
	}
	return a
}

type App struct {
	*Server
	*echo.Group  //没有指定域名时有效
	http.Handler //指定域名时有效
	Name         string
	Domain       string
	controllers  map[string]*Controller
}

func (a *App) G() *echo.Group {
	return a.Group
}

func (a *App) E() *echo.Echo {
	return a.Handler.(*echo.Echo)
}

//注册路由：app.R(`/index`,Index.Index,"GET","POST")
func (a *App) R(path string, h echo.Handler, methods ...string) *App {
	if len(methods) < 1 {
		methods = append(methods, "GET")
	}
	a.Webx().Match(methods, path, h)
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
func (a *App) RC(c interface{}, args ...echo.HandlerFunc) *Controller {
	name := fmt.Sprintf("%T", c)
	cr := &Controller{
		Controller: c,
		Webx:       a.Webx(),
	}
	switch len(args) {
	case 1:
		cr.Before = args[0]
	case 2:
		cr.Before = args[0]
		cr.After = args[1]
	}
	a.controllers[name] = cr
	return cr
}
