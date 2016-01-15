package xsrf

import (
	"errors"
	"html/template"
	"net/http"
	"time"

	"github.com/webx-top/webx/echo"
	"github.com/webx-top/webx/lib/uuid"
)

func NewXsrf(args ...Manager) *Xsrf {
	x := &Xsrf{
		FieldName: `_xsrf`,
		On:        true,
	}
	if len(args) > 0 {
		x.Manager = args[0]
	} else {
		x.Manager = &CookieStorage{
			Secret: time.Now().Local().String() + `@webx.top`,
			Codec:  codec.Default,
		}
	}
	return x
}

type Xsrf struct {
	Manager
	FieldName string
	On        bool
}

func (c *Xsrf) Value(ctx *echo.Context) string {
	var val string = c.Manager.Get(c.FieldName, ctx)
	if val == "" {
		val = uuid.NewRandom().String()
		c.Manager.Set(c.FieldName, val, ctx)
	}
	return val
}

func (c *Xsrf) Form(ctx *echo.Context) template.HTML {
	var html string
	if c.On {
		html = fmt.Sprintf(`<input type="hidden" name="%v" value="%v" />`, c.FieldName, c.Value(ctx))
	}
	return template.HTML(html)
}

func (c *Xsrf) Middleware() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			if !c.On {
				return h(ctx)
			}
			if ignore, _ := ctx.Get(`webx:ignoreXsrf`).(bool); ignore {
				return h(ctx)
			}
			val := c.Value(ctx)
			if ctx.Request().Method == `POST` {
				formVal := ctx.Form(c.FieldName)
				if formVal == "" || val != formVal {
					return errors.New("xsrf token error.")
				}
			}
			return h(ctx)
		}
	}
}

type Manager interface {
	Get(key string, ctx *echo.Context) string
	Set(key, val string, ctx *echo.Context)
	Valid(key, val string, ctx *echo.Context) bool
}
