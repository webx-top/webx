package webx

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func NewApp(name string, domain string, s *Server, middlewares ...echo.Middleware) (a *App) {
	a = &App{
		Server:      s,
		Name:        name,
		Domain:      domain,
		controllers: make(map[string]interface{}),
	}
	if a.Domain == "" {
		var prefix string
		if name != "" {
			prefix = `/` + name
		}
		a.Group = s.Echo.Group(prefix, middlewares...)
	} else {
		e := echo.New()
		e.Use(s.DefaultMiddlewares...)
		e.Use(middlewares...)
		if s.TemplateEngine != nil {
			e.SetRenderer(s.TemplateEngine)
		}
		a.Handler = e
	}
	return
}

type App struct {
	*Server
	*echo.Group  //没有指定域名时有效
	http.Handler //指定域名时有效
	Name         string
	Domain       string
	controllers  map[string]interface{}
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
	if a.Group != nil {
		a.G().Match(methods, path, h)
	} else {
		a.E().Match(methods, path, h)
	}
	return a
}

//获取控制器
func (a *App) C(name string) (c interface{}) {
	c, _ = a.controllers[name]
	return
}

//登记控制器
func (a *App) RC(c interface{}) *App {
	name := fmt.Sprintf("%T", c)
	a.controllers[name] = c
	return a
}
