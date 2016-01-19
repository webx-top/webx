package htmlcache

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/webx-top/echo"
)

func OutputXML(content []byte, ctx echo.Context, args ...int) (err error) {
	var code int = http.StatusOK
	if len(args) > 0 {
		code = args[0]
	}
	ctx.X().Xml(code, content)
	return nil
}

func OutputJSON(content []byte, ctx echo.Context, args ...int) (err error) {
	callback := ctx.Query(`callback`)
	var code int = http.StatusOK
	if len(args) > 0 {
		code = args[0]
	}
	if callback != `` {
		ctx.X().Jsonp(code, callback, content)
	} else {
		ctx.X().Json(code, content)
	}
	return nil
}

func OutputHTML(content []byte, ctx echo.Context, args ...int) (err error) {
	var code int = http.StatusOK
	if len(args) > 0 {
		code = args[0]
	}
	return ctx.HTML(code, string(content))
}

func RenderXML(ctx echo.Context) (b []byte, err error) {
	b, err = xml.Marshal(ctx.Get(`Data`))
	return
}

func RenderJSON(ctx echo.Context) (b []byte, err error) {
	b, err = json.Marshal(ctx.Get(`Data`))
	return
}

func RenderHTML(ctx echo.Context) (b []byte, err error) {
	tmpl, _ := ctx.Get(`Tmpl`).(string)
	if tmpl == `` {
		return
	}
	b, err = ctx.X().Fetch(tmpl, ctx.Get(`Data`))
	return
}

func Output(format string, content []byte, ctx echo.Context) (err error) {
	switch format {
	case `xml`:
		return OutputXML(content, ctx)
	case `json`:
		return OutputJSON(content, ctx)
	default:
		return OutputHTML(content, ctx)
	}
}

func Render(ctx echo.Context, args ...int) error {
	format, ok := ctx.Get(`webx:format`).(string)
	if !ok {
		format = ctx.Query(`format`)
	}
	switch format {
	case `xml`:
		b, err := RenderXML(ctx)
		if err != nil {
			return err
		}
		return OutputXML(b, ctx, args...)
	case `json`:
		b, err := RenderJSON(ctx)
		if err != nil {
			return err
		}
		return OutputJSON(b, ctx, args...)
	default:
		tmpl, _ := ctx.Get(`Tmpl`).(string)
		if tmpl == `` {
			return nil
		}
		var code int = http.StatusOK
		if len(args) > 0 {
			code = args[0]
		}
		return ctx.Render(code, tmpl, ctx.Get(`Data`))
	}
}
