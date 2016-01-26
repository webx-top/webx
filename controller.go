package webx

import (
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

func (a *Controller) Init(c *Context) {
	a.Context = c
}

func (a *Controller) Before() error {
	a.SetFunc("Query", a.Query)
	a.SetFunc("Form", a.Form)
	a.SetFunc("Path", a.Path)
	return nil
}

func (a *Controller) After() error {
	return a.Display()
}

func (a *Controller) X(c echo.Context) *Context {
	return X(c)
}
