package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/webx-top/echo"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/htmlcache"
	"github.com/webx-top/webx/lib/middleware/language"
	"github.com/webx-top/webx/lib/middleware/session"
)

type Index struct {
	*X.App
}

func (i *Index) Index(c echo.Context) error {
	fmt.Println(`Index.`)
	c.Set(`webx:tmpl`, `index`)
	return nil
}

func (i *Index) Index2(c echo.Context) error {
	return c.Render(http.StatusOK, `index2`, nil)
}

func (i *Index) Before(c echo.Context) error {
	fmt.Println(`Before.`)
	if Cfg.Read(c) {
		c.X().Echo().Logger().Info(`htmlcache valid.`)
		return nil
	}
	c.X().Echo().Logger().Info(`htmlcache invalid.`)
	return nil
}

func (i *Index) After(c echo.Context) error {
	fmt.Println(`After.`)

	//=========================================
	tmpl := X.MustString(c, `webx:tmpl`)
	if tmpl != `` {
		buf := new(bytes.Buffer)
		if err := i.App.Server.TemplateEngine.Render(buf, tmpl, c.Get(`Data`), nil); err != nil {
			return err
		}
		if Cfg.Write(buf.Bytes(), c) {
			c.X().Echo().Logger().Info(`htmlcache wroten.`)
		}
		return c.HTML(200, buf.String())
	}
	return nil
}

var indexController *Index
var Cfg = &htmlcache.Config{
	HtmlCacheDir:   `html`,
	HtmlCacheOn:    true,
	HtmlCacheRules: make(map[string]interface{}),
	HtmlCacheTime:  86400,
}

func main() {
	var lang = language.NewLanguage()
	lang.Set(`zh-cn`, true, true)
	lang.Set(`en`, true)
	var store = session.NewCookieStore([]byte("secret-key"))

	s := X.Serv().InitTmpl().Pprof().Debug(true).SetHook(lang.DetectURI)
	Cfg.HtmlCacheRules[`index:`] = []interface{}{
		`index.html`, /*/保存名称
		func(tmpl string, c echo.Context) string { //自定义保存名称
			return tmpl
		},
		func(tmpl string, c echo.Context) (mtime int64,expired bool) { //判断缓存是否过期
			return
		},*/
	}
	Cfg.HtmlCacheRules[`test:`] = []interface{}{
		`test.html`,
	}

	//==================================
	//测试多语言切换和session
	//==================================
	app := s.NewApp("", lang.Store(), session.Sessions("XSESSION", store))
	indexController = &Index{App: app}
	//测试session
	app.R("/", func(c echo.Context) error {

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

		return c.String(http.StatusOK, fmt.Sprintf(`Hello world.Count:%v.Language: %v`, count, c.Get(language.LANG_KEY)))
	}).
		R("/t", func(c echo.Context) error {
		return c.Render(http.StatusOK, `index`, nil)
	}, `GET`).
		//测试Before和After以及全页面html缓存
		RC(indexController, indexController.Before, indexController.After).
		R("/index", indexController.Index).
		R("/index2", indexController.Index2)

	//=======================================
	//测试以中间件形式实现的全页面缓存功能
	//=======================================
	s.NewApp("test", Cfg.Middleware()).
		R("", func(c echo.Context) error {
		c.Set(`Tmpl`, `index2`)
		return nil
	}, `GET`)

	//=======================================
	//测试无任何中间件时是否正常
	//=======================================
	s.NewApp("ping").R("", func(c echo.Context) error {
		return c.String(200, "pong")
	})

	s.Run("127.0.0.1", "8080")
}
