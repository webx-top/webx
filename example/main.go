package main

import (
	"fmt"
	"net/http"

	X "bitbucket.org/admpub/webx"
	MW "bitbucket.org/admpub/webx/lib/middleware"
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
	var lang = MW.NewLanguage()
	lang.Set(`zh-cn`, true, true)
	lang.Set(`en`, true)

	s := X.Serv().InitTmpl().Pprof().Debug(true).SetHook(lang.DetectURI)
	s.NewApp("", lang.Store()).R("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, `Hello world.Language: `+fmt.Sprintf("%v", c.Get(`language`)))
	}).
		R("/t", func(c *echo.Context) error {
		return c.Render(http.StatusOK, `index`, nil)
	}, `GET`).
		RC(indexController).
		R("/index", indexController.Index).
		R("/index2", indexController.Index2)

	s.Run("127.0.0.1", "8080")
}
