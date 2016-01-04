package webx

import (
	"net/http"

	"github.com/labstack/echo"
)

func NewApp(name string, e *echo.Echo) (a *App) {
	a = &App{
		Name:    name,
		Handler: e,
	}
	return
}

type App struct {
	Name string
	http.Handler
}
