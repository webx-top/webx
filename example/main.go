package main

import (
	"net/http"

	X "bitbucket.org/admpub/webx"
	"github.com/admpub/echo"
)

type Index struct {
}

func (i *Index) Index(c *echo.Context) error {
	return c.Render(http.StatusOK, `index`, nil)
}

func (i *Index) Index2(c *echo.Context) error {
	return c.Render(http.StatusOK, `index2`, nil)
}

var indexController *Index = &Index{}

func main() {
	s := X.Serv().InitTmpl().Pprof().Debug(true)
	s.NewApp("").R("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, `Hello world.`)
	}).
		R("/t", func(c *echo.Context) error {
			return c.Render(http.StatusOK, `index`, nil)
		}, `GET`).
		RC(indexController).
		R("/index", indexController.Index).
		R("/index2", indexController.Index2)

	s.Run("127.0.0.1", "8080")
}
