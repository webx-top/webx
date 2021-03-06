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
package xsrf

import (
	"errors"
	"fmt"
	"html/template"

	"github.com/webx-top/echo"
	"github.com/webx-top/webx/lib/com"
	"github.com/webx-top/webx/lib/uuid"
)

func New(args ...Manager) *Xsrf {
	x := &Xsrf{
		FieldName: `_xsrf`,
		On:        true,
	}
	if len(args) > 0 {
		x.Manager = args[0]
	} else {
		x.Manager = &SecCookieStorage{}
		//x.Manager = &CookieStorage{}
		//x.Manager = &SessionStorage{}
	}
	return x
}

type Xsrf struct {
	Manager
	FieldName string
	On        bool
}

func (c *Xsrf) Value(ctx echo.Context) string {
	var val string = c.Manager.Get(c.FieldName, ctx)
	if val == "" {
		val = uuid.NewRandom().String()
		c.Manager.Set(c.FieldName, val, ctx)
	}
	return val
}

func (c *Xsrf) Form(ctx echo.Context) template.HTML {
	var html string
	if c.On {
		html = fmt.Sprintf(`<input type="hidden" name="%v" value="%v" />`, c.FieldName, com.HtmlEncode(c.Value(ctx)))
	}
	return template.HTML(html)
}

func (c *Xsrf) Ignore(on bool, ctx echo.Context) {
	ctx.Set(`webx:ignoreXsrf`, on)
}

func (c *Xsrf) Register(ctx echo.Context) {
	ctx.SetFunc("XsrfForm", func() template.HTML {
		return c.Form(ctx)
	})
	ctx.SetFunc("XsrfValue", func() string {
		return c.Value(ctx)
	})
	ctx.SetFunc("XsrfName", func() string {
		return c.FieldName
	})
}

func (c *Xsrf) Middleware() echo.MiddlewareFunc {
	return echo.MiddlewareFunc(func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(ctx echo.Context) error {
			if !c.On {
				return h.Handle(ctx)
			}
			if ignore, _ := ctx.Get(`webx:ignoreXsrf`).(bool); ignore {
				return h.Handle(ctx)
			}
			c.Register(ctx)
			val := c.Value(ctx)
			if ctx.Request().Method() == `POST` {
				formVal := ctx.Form(c.FieldName)
				if formVal == "" || val != formVal {
					return errors.New("xsrf token error.")
				}
			}
			return h.Handle(ctx)
		})
	})
}

type Manager interface {
	Get(key string, ctx echo.Context) string
	Set(key, val string, ctx echo.Context)
	Valid(key, val string, ctx echo.Context) bool
}
