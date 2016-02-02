package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/webx-top/echo"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/htmlcache"
	"github.com/webx-top/webx/lib/middleware/language"
)

type Index struct {
	index  X.Mapper
	index2 X.Mapper
	*X.Controller
}

func (a *Index) Init(c *X.Context) {
	a.Controller = X.NewController(c)
}

func (a *Index) Index() error {
	fmt.Println(`Index.`)
	a.Tmpl = `index`
	return nil
}

func (a *Index) Index2() error {
	a.Tmpl = `index2`
	return nil
}

func (a *Index) Before() error {
	fmt.Println(`Before.`)
	return a.Controller.Before()
}

func (a *Index) After() error {
	fmt.Println(`After.`)
	return a.Controller.After()
}

var indexController *Index
var Cfg = &htmlcache.Config{
	HtmlCacheDir:   `html`,
	HtmlCacheOn:    true,
	HtmlCacheRules: make(map[string]interface{}),
	HtmlCacheTime:  86400,
}

func main() {
	mode := flag.String("m", "clean", "port of your app.")
	flag.Parse()

	var lang = language.NewLanguage()
	lang.Set(`zh-cn`, true, true)
	lang.Set(`en`, true)

	Cfg.HtmlCacheRules[`index:`] = []interface{}{
		`index`, /*/保存名称
		func(tmpl string, c echo.Context) string { //自定义保存名称
			return tmpl
		},
		func(tmpl string, c echo.Context) (mtime int64,expired bool) { //判断缓存是否过期
			return
		},*/
	}
	Cfg.HtmlCacheRules[`test:`] = []interface{}{
		`test`,
	}
	var s *X.Server
	if *mode == `clean` {
		// ===============================================================
		// benchmark测试(不使用任何中间件，特别是log中间件，比较影响速度)
		// ===============================================================
		s = X.Serv()
		s.DefaultMiddlewares = []echo.Middleware{}
		s.Core = echo.New(s.InitContext)
	} else {
		s = X.Serv().InitTmpl().Pprof().Debug(true).SetHook(lang.DetectURI)
	}

	//==================================
	//测试多语言切换和session
	//==================================
	app := s.NewApp("", lang.Store(), Cfg.Middleware())
	indexController = &Index{}
	//测试session
	app.R("/", func(c *X.Context) error {
		var count int
		v := c.GetSession("count")

		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count += 1
		}

		c.SetSession("count", count)

		return c.String(http.StatusOK, fmt.Sprintf(`Hello world.Count:%v.Language: %v`, count, c.Language))
	}).
		R("/t", func(c *X.Context) error {
			return c.Render(http.StatusOK, `index`, nil)
		}, `GET`).
		//测试Before和After以及全页面html缓存
		Reg(indexController).Auto()

	//=======================================
	//测试以中间件形式实现的全页面缓存功能
	//=======================================
	s.NewApp("test", Cfg.Middleware()).
		R("", func(c *X.Context) error {
			c.Tmpl = `index2`
			return nil
		}, `GET`)

	//=======================================
	//测试无任何中间件时是否正常
	//=======================================
	s.NewApp("ping").R("", func(c *X.Context) error {
		return c.String(200, "pong")
	})

	s.Run("127.0.0.1", "8080")
}
