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
	"fmt"
	"strings"

	"github.com/webx-top/echo"
)

type Webxer interface {
	URL(echo.Handler, ...interface{}) string
	SetRenderer(echo.Renderer)
	Use(...echo.Middleware)
	PreUse(...echo.Middleware)
	Any(string, echo.Handler, ...echo.Middleware)
	Match([]string, string, echo.Handler, ...echo.Middleware)
	Trace(string, echo.Handler, ...echo.Middleware)
	Group(string, ...echo.Middleware) *echo.Group
}

func NewApp(name string, domain string, s *Server, middlewares ...echo.Middleware) (a *App) {
	a = &App{
		Server:      s,
		Name:        name,
		Domain:      domain,
		controllers: make(map[string]*Wrapper),
	}
	if s.TemplateEngine != nil {
		a.Renderer = s.TemplateEngine
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
		a.Group = s.Core.Group(prefix)
		a.Group.Use(middlewares...)
	} else {
		e := echo.NewWithContext(s.InitContext)
		e.Use(s.DefaultMiddlewares...)
		e.Use(middlewares...)
		a.Handler = e
		a.Url = `http://` + a.Domain + `/`
		a.Dir = `/`
	}
	return
}

type App struct {
	*Server
	*echo.Group            //没有指定域名时有效
	Handler     *echo.Echo //指定域名时有效
	Renderer    echo.Renderer
	Name        string
	Domain      string
	controllers map[string]*Wrapper
	Url         string
	Dir         string
}

func (a *App) G() *echo.Group {
	return a.Group
}

func (a *App) E() *echo.Echo {
	return a.Handler
}

// 注册路由：app.R(`/index`,Index.Index,"GET","POST")
func (a *App) R(path string, h HandlerFunc, methods ...string) *App {
	if len(methods) < 1 {
		methods = append(methods, "GET")
	}
	_, ctl, act := a.Server.URL.Set(path, h)
	a.Webx().Match(methods, path, echo.HandlerFunc(func(ctx echo.Context) error {
		c := X(ctx)
		if err := c.Init(a, nil, ctl, act); err != nil {
			return err
		}
		return h(c)
	}))
	return a
}

func (a *App) Webx() Webxer {
	if a.Group != nil {
		return a.G()
	}
	return a.E()
}

// 获取控制器
func (a *App) Ctl(name string) (c interface{}) {
	c, _ = a.controllers[name]
	return
}

// 登记控制器
func (a *App) Reg(c interface{}) *Wrapper {
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
