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
package htmlcache

import (
	"encoding/json"
	"encoding/xml"
	"net/http"

	X "github.com/webx-top/webx"
)

func OutputXML(content []byte, ctx *X.Context, args ...int) (err error) {
	code := ctx.Code
	if len(args) > 0 {
		code = args[0]
	}
	ctx.Object().XMLBlob(code, content)
	return nil
}

func OutputJSON(content []byte, ctx *X.Context, args ...int) (err error) {
	callback := ctx.Query(`callback`)
	code := ctx.Code
	if len(args) > 0 {
		code = args[0]
	}
	if callback != `` {
		content = []byte(callback + "(" + string(content) + ");")
	}
	ctx.Object().JSONBlob(code, content)
	return nil
}

func OutputHTML(content []byte, ctx *X.Context, args ...int) (err error) {
	code := ctx.Code
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
	ctx.Context.SetFunc(`Status`, func() int {
		return ctx.Output.Status
	})
	ctx.Context.SetFunc(`Message`, func() interface{} {
		return ctx.Output.Message
	})
	b, err = ctx.Object().Fetch(ctx.Tmpl, ctx.Output.Data)
	return
}

func Output(content []byte, ctx *X.Context) (err error) {
	ctx.Code = http.StatusOK
	switch ctx.Format {
	case `xml`:
		return OutputXML(content, ctx)
	case `json`:
		return OutputJSON(content, ctx)
	default:
		return OutputHTML(content, ctx)
	}
}
