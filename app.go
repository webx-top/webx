package webx

import (
	"fmt"
	"net/http"

	"bitbucket.org/admpub/webx/lib/pprof"
	"github.com/labstack/echo"
)

func NewApp(name string, e *echo.Echo) (a *App) {
	a = &App{
		Name:        name,
		Handler:     e,
		controllers: make(map[string]interface{}),
	}
	return
}

type App struct {
	Name string
	http.Handler
	controllers map[string]interface{}
}

func (a *App) E() *echo.Echo {
	return a.Handler.(*echo.Echo)
}

//注册路由：app.R(`/index`,Index.Index,"GET","POST")
func (a *App) R(path string, h echo.Handler, methods ...string) *App {
	if len(methods) < 1 {
		methods = append(methods, "GET")
	}
	a.E().Match(methods, path, h)
	return a
}

//创建新Group路由
func (a *App) NewG(prefix string, m ...echo.Middleware) *echo.Group {
	return a.E().Group(prefix, m...)
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

//启用pprof
func (a *App) Pprof() *App {
	pprof.Wrapper(a.E())
	return a
}

//开关debug模式
func (a *App) Debug(on bool) *App {
	a.E().SetDebug(on)
	return a
}
