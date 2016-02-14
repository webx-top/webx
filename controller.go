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
	"net/http"

	"github.com/webx-top/echo"
)

func NewController(c *Context) *Controller {
	a := &Controller{}
	a.Init(c)
	return a
}

type Controller struct {
	*Context
}

func (a *Controller) Init(c *Context) error {
	a.Context = c
	a.SetFunc("Query", a.Query)
	a.SetFunc("Form", a.Form)
	a.SetFunc("Path", a.Path)
	return nil
}

func (a *Controller) X(c echo.Context) *Context {
	return X(c)
}

func (a *Controller) Redirect(url string, args ...interface{}) error {
	var code = http.StatusFound //302. 307:http.StatusTemporaryRedirect
	if len(args) > 0 {
		if v, ok := args[0].(bool); ok && v {
			code = http.StatusMovedPermanently
		} else if v, ok := args[0].(int); ok {
			code = v
		}
	}
	a.Context.Exit = true
	if a.Format != `html` {
		a.Context.Set(`webx:ignoreRender`, false)
		a.Assign(`Location`, url)
		return a.Display()
	}
	return a.Context.Redirect(code, url)
}

func (a *Controller) NotFound(args ...string) error {
	var code = http.StatusNotFound
	var text = "Page not found"
	a.Context.Exit = true
	if len(args) > 0 {
		text = args[0]
	}
	return a.Errno(code, text)
}
