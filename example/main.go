package main

import (
	"net/http"

	X "bitbucket.org/admpub/webx"
	"github.com/labstack/echo"
)

type Index struct {
}

func (i *Index) Index(c *echo.Context) error {
	return c.Render(http.StatusOK, `index`, nil)
}

var indexController *Index = &Index{}

func main() {
	s := X.Serv()
	s.InitTmpl().NewApp("*").Pprof().Debug(true).
		R("/", func(c *echo.Context) error {
			return c.String(http.StatusOK, `Hello world.`)
		}).
		R("/t", func(c *echo.Context) error {
			return c.Render(http.StatusOK, `index`, nil)
		}, `GET`).
		RC(indexController).
		R("/index", indexController.Index)

	s.Run("127.0.0.1", "8080")
}
