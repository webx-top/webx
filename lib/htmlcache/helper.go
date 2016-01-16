package htmlcache

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/webx-top/echo"
)

func OutputXML(content []byte, ctx *echo.Context, args ...int) (err error) {
	ctx.Response().Header().Set(echo.ContentType, echo.ApplicationXMLCharsetUTF8)
	var code int = http.StatusOK
	if len(args) > 0 {
		code = args[0]
	}
	ctx.Response().WriteHeader(code)
	ctx.Response().Write([]byte(xml.Header))
	ctx.Response().Write(content)
	return nil
}

func OutputJSON(content []byte, ctx *echo.Context, args ...int) (err error) {
	callback := ctx.Query(`callback`)
	ctx.Response().Header().Set(echo.ContentType, echo.ApplicationJSONCharsetUTF8)
	var code int = http.StatusOK
	if len(args) > 0 {
		code = args[0]
	}
	ctx.Response().WriteHeader(code)
	if callback != `` {
		ctx.Response().Write([]byte(callback + "("))
		ctx.Response().Write(content)
		ctx.Response().Write([]byte(");"))
	} else {
		ctx.Response().Write(content)
	}
	return nil
}

func OutputHTML(content []byte, ctx *echo.Context, args ...int) (err error) {
	var code int = http.StatusOK
	if len(args) > 0 {
		code = args[0]
	}
	return ctx.HTML(code, string(content))
}

func RenderXML(ctx *echo.Context) (b []byte, err error) {
	b, err = xml.Marshal(ctx.Get(`Data`))
	return
}

func RenderJSON(ctx *echo.Context) (b []byte, err error) {
	b, err = json.Marshal(ctx.Get(`Data`))
	return
}

func RenderHTML(renderer echo.Renderer, ctx *echo.Context) (b []byte, err error) {
	tmpl, _ := ctx.Get(`Tmpl`).(string)
	if tmpl == `` {
		return
	}
	buf := new(bytes.Buffer)
	err = renderer.Render(buf, tmpl, ctx.Get(`Data`), ctx.Funcs)
	b = buf.Bytes()
	return
}

func Output(format string, content []byte, ctx *echo.Context) (err error) {
	switch format {
	case `xml`:
		return OutputXML(content, ctx)
	case `json`:
		return OutputJSON(content, ctx)
	default:
		return OutputHTML(content, ctx)
	}
}

func Render(ctx *echo.Context, args ...int) error {
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
