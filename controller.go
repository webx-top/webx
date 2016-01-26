package webx

import (
	"github.com/webx-top/echo"
)

func NewController() *Controller {
	return &Controller{}
}

type Controller struct {
	*Context
	*App
}

func (a *Controller) Init(c *Context, app *App) {
	a.Context = c
	a.App = app
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
