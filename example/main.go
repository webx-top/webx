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

func (a *Index) Init(c *X.Context) error {
	a.Controller = X.NewController(c)
	return nil
}

func (a *Index) Index() error {
	fmt.Println(`Index.`)
	return a.Display(`index`)
}

func (a *Index) Index2() error {
	return a.Display(`index2`)
}

func (a *Index) Before() error {
	fmt.Println(`Before.`)
	return nil
}

func (a *Index) After() error {
	fmt.Println(`After.`)
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
	mode := flag.String("m", "clean2", "port of your app.")
	port := flag.String("p", "8080", "port of your app.")
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
		//s.ResetTmpl().Pprof().Debug(true)
		s.DefaultMiddlewares = []echo.Middleware{}
		s.Core = echo.NewWithContext(s.InitContext)
	} else {
		s = X.Serv().ResetTmpl().Pprof().Debug(true)
	}

	//==================================
	//测试多语言切换和session
	//==================================
	app := s.NewApp("", Cfg.Middleware())
	s.Core.PreUse(lang.Middleware())
	indexController = &Index{}
	//测试session
	app.R("/session", func(c *X.Context) error {
		var count int
		v := c.Session().Get("count")

		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count += 1
		}

		c.Session().Set("count", count).Save()

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
			return c.Display(`index2`)
		}, `GET`)

	//=======================================
	//测试无任何中间件时是否正常
	//=======================================
	s.NewApp("ping").R("", func(c *X.Context) error {
		return c.String(200, "pong")
	})

	s.Run(":" + *port)
}
