package main

import (
	"fmt"
	"net/http"

	X "bitbucket.org/admpub/webx"
	MW "bitbucket.org/admpub/webx/lib/middleware"
	"bitbucket.org/admpub/webx/lib/middleware/session"
	"github.com/admpub/echo"
)

type Index struct {
}

func (i *Index) Index(c *echo.Context) error {
	fmt.Println(`Index.`)
	return c.Render(http.StatusOK, `index`, nil)
}

func (i *Index) Index2(c *echo.Context) error {
	return c.Render(http.StatusOK, `index2`, nil)
}

func (i *Index) Before(c *echo.Context) error {
	fmt.Println(`Before.`)
	return nil
}

func (i *Index) After(c *echo.Context) error {
	fmt.Println(`After.`)
	return nil
}

var indexController *Index = &Index{}

func main() {
	var lang = MW.NewLanguage()
	lang.Set(`zh-cn`, true, true)
	lang.Set(`en`, true)
	var store = session.NewCookieStore([]byte("secret-key"))

	s := X.Serv().InitTmpl().Pprof().Debug(true).SetHook(lang.DetectURI)
	s.NewApp("", lang.Store(), session.Sessions("XSESSION", store)).R("/", func(c *echo.Context) error {

		session := session.Default(c)
		var count int
		v := session.Get("count")

		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count += 1
		}

		session.Set("count", count)
		session.Save()

		return c.String(http.StatusOK, fmt.Sprintf(`Hello world.Count:%v.Language: %v`, count, c.Get(`language`)))
	}).
		R("/t", func(c *echo.Context) error {
		return c.Render(http.StatusOK, `index`, nil)
	}, `GET`).
		RC(indexController, indexController.Before, indexController.After).
		R("/index", indexController.Index).
		R("/index2", indexController.Index2)

	s.Run("127.0.0.1", "8080")
}
