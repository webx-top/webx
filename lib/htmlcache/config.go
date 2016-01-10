package htmlcache

import (
	"bytes"
	"net/http"
	"strings"
	"time"

	X "bitbucket.org/admpub/webx"
	"bitbucket.org/admpub/webx/lib/com"
	"github.com/admpub/echo"
)

type Config struct {
	HtmlCacheDir   string
	HtmlCacheOn    bool
	HtmlCacheRules map[string]interface{}
	HtmlCacheTime  interface{}
}

func (c *Config) Read(ctx *echo.Context) bool {
	req := ctx.Request()
	if !c.HtmlCacheOn || req.Method != `GET` {
		return false
	}
	p := strings.Trim(req.URL.Path, `/`)
	if p == `` {
		p = `index`
	}
	s := strings.SplitN(p, `/`, 3)
	var rule *Rule
	switch len(s) {
	case 2:
		k := s[0] + `:` + s[1]
		if v, ok := c.HtmlCacheRules[k]; ok {
			rule = c.Rule(v)
		} else if v, ok := c.HtmlCacheRules[s[1]]; ok {
			rule = c.Rule(v)
		} else {
			k = s[0] + `:`
			if v, ok := c.HtmlCacheRules[k]; ok {
				rule = c.Rule(v)
			}
		}
	case 1:
		k := s[0] + `:`
		if v, ok := c.HtmlCacheRules[k]; ok {
			rule = c.Rule(v)
		}
	}
	var saveFile string = c.SaveFileName(rule, ctx)
	if saveFile == "" {
		return false
	}
	mtime, expired := c.Expired(rule, ctx, saveFile)
	if expired {
		ctx.Set(`webx:saveHtmlFile`, saveFile)
		return false
	}
	if !HttpCache(ctx, mtime, nil) {
		html, err := com.ReadFileS(saveFile)
		if err != nil {
			ctx.Echo().Logger().Error(err)
		}
		ctx.HTML(http.StatusOK, html)
	}
	ctx.Set(`webx:exit`, true)
	return true
}

func (c *Config) Rule(rule interface{}) *Rule {
	r := &Rule{}
	switch rule.(type) {
	case Rule:
		v := rule.(Rule)
		r = &v
	case *Rule:
		r = rule.(*Rule)
	case []interface{}:
		v := rule.([]interface{})
		switch len(v) {
		case 3:
			switch v[2].(type) {
			case int:
				r.ExpireTime = v[2].(int)
			case func(string, *echo.Context) (int64, bool):
				r.ExpireFunc = v[2].(func(string, *echo.Context) (int64, bool))
			}
			fallthrough
		case 2:
			r.SaveFunc = v[1].(func(string, *echo.Context) string)
			fallthrough
		case 1:
			r.SaveFile = v[0].(string)
		default:
			return nil
		}
	case string:
		r.SaveFile = rule.(string)
	default:
		return nil
	}
	return r
}

func (c *Config) Write(buf *bytes.Buffer, ctx *echo.Context) bool {
	if !c.HtmlCacheOn || ctx.Request().Method != `GET` {
		return false
	}
	tmpl := X.MustString(ctx, `webx:saveHtmlFile`)
	if tmpl == `` {
		return false
	}
	if err := com.WriteFile(tmpl, buf.Bytes()); err != nil {
		ctx.Echo().Logger().Debug(err)
	}
	return true
}

func (c *Config) SaveFileName(rule *Rule, ctx *echo.Context) string {
	if rule == nil {
		return ""
	}
	var saveFile string = rule.SaveFile
	if rule.SaveFunc != nil {
		saveFile = rule.SaveFunc(saveFile, ctx)
	}
	return saveFile
}

func (c *Config) Expired(rule *Rule, ctx *echo.Context, saveFile string) (int64, bool) {
	var expired int64
	if rule.ExpireTime > 0 {
		expired = int64(rule.ExpireTime)
	} else if rule.ExpireFunc != nil {
		return rule.ExpireFunc(saveFile, ctx)
	} else {
		switch c.HtmlCacheTime.(type) {
		case int:
			expired = int64(c.HtmlCacheTime.(int))
		case int64:
			expired = c.HtmlCacheTime.(int64)
		case func(string, *echo.Context) (int64, bool):
			fn := c.HtmlCacheTime.(func(string, *echo.Context) (int64, bool))
			return fn(saveFile, ctx)
		}
	}
	mtime, err := com.FileMTime(saveFile)
	if err != nil {
		ctx.Echo().Logger().Debug(err)
	}
	if mtime == 0 {
		return mtime, true
	}
	if time.Now().Local().Unix() > mtime+expired {
		return mtime, true
	}
	return mtime, false
}

func (c *Config) Middleware(renderer echo.Renderer) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx *echo.Context) error {
			if c.Read(ctx) {
				return nil
			}
			if err := h(ctx); err != nil {
				return err
			}
			tmpl, _ := ctx.Get(`Tmpl`).(string)
			if tmpl == `` {
				return nil
			}
			buf := new(bytes.Buffer)
			if err := renderer.Render(buf, tmpl, ctx.Get(`Data`)); err != nil {
				return err
			}
			c.Write(buf, ctx)
			return ctx.HTML(http.StatusOK, buf.String())
		}
	}
}
