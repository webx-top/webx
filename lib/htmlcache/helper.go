package htmlcache

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/webx-top/echo"
	X "github.com/webx-top/webx"
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

func RenderXML(ctx *X.Context) (b []byte, err error) {
	b, err = xml.Marshal(ctx.Output)
	return
}

func RenderJSON(ctx *X.Context) (b []byte, err error) {
	b, err = json.Marshal(ctx.Output)
	return
}

func RenderHTML(ctx *X.Context) (b []byte, err error) {
	if ctx.Tmpl == `` {
		return
	}
	b, err = ctx.X().Fetch(ctx.Tmpl, ctx.Output)
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
