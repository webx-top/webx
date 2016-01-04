package main

import (
	"net/http"

	X "bitbucket.org/admpub/webx"
	"bitbucket.org/admpub/webx/lib/pprof"
	"github.com/labstack/echo"
)

func main() {
	s := X.Serv().Template()
	e := s.New("*")
	e.Debug()
	e.Get("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, `Hello world.`)
	})
	e.Get("/t", func(c *echo.Context) error {
		return c.Render(http.StatusOK, `index`, nil)
	})
	pprof.Wrapper(e)
	s.Run("127.0.0.1", "8080")
}
