package htmlcache

import (
	"strings"
	"time"

	"github.com/webx-top/echo"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/com"
)

type Config struct {
	HtmlCacheDir   string
	HtmlCacheOn    bool
	HtmlCacheRules map[string]interface{}
	HtmlCacheTime  interface{}
	htmlCacheRules map[string]*Rule
}

func (c *Config) Read(ctx echo.Context) bool {
	format := ctx.Query(`format`)
	ctx.Set(`webx:format`, format)
	req := ctx.Request()
	if !c.HtmlCacheOn || req.Method != `GET` {
		return false
	}
	p := strings.Trim(req.URL.Path, `/`)
	if p == `` {
		p = `index`
	}
	s := strings.SplitN(p, `/`, 3)

	if c.htmlCacheRules == nil {
		c.htmlCacheRules = make(map[string]*Rule)
		for key, rule := range c.HtmlCacheRules {
			c.htmlCacheRules[key] = c.Rule(rule)
		}
	}

	var rule *Rule
	switch len(s) {
	case 2:
		k := s[0] + `:` + s[1]
		if v, ok := c.htmlCacheRules[k]; ok {
			rule = v
		} else if v, ok := c.htmlCacheRules[s[1]]; ok {
			rule = v
		} else {
			k = s[0] + `:`
			if v, ok := c.htmlCacheRules[k]; ok {
				rule = v
			}
		}
	case 1:
		k := s[0] + `:`
		if v, ok := c.htmlCacheRules[k]; ok {
			rule = v
		}
	}
	var saveFile string = c.SaveFileName(rule, ctx)
	if saveFile == "" {
		return false
	}
	if format != `` {
		saveFile += `.` + format
	}
	mtime, expired := c.Expired(rule, ctx, saveFile)
	if expired {
		ctx.Set(`webx:saveHtmlFile`, saveFile)
		return false
	}
	//_ = mtime
	//ctx.File(saveFile, ``, false)

	if !HttpCache(ctx, mtime, nil) {
		html, err := com.ReadFile(saveFile)
		if err != nil {
			ctx.X().Echo().Logger().Error(err)
		}
		Output(format, html, ctx)
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
			case func(string, echo.Context) (int64, bool):
				r.ExpireFunc = v[2].(func(string, echo.Context) (int64, bool))
			}
			fallthrough
		case 2:
			r.SaveFunc = v[1].(func(string, echo.Context) string)
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

func (c *Config) Write(b []byte, ctx echo.Context) bool {
	if !c.HtmlCacheOn || ctx.Request().Method != `GET` {
		return false
	}
	tmpl := X.MustString(ctx, `webx:saveHtmlFile`)
	if tmpl == `` {
		return false
	}
	if err := com.WriteFile(tmpl, b); err != nil {
		ctx.X().Echo().Logger().Debug(err)
	}
	return true
}

func (c *Config) SaveFileName(rule *Rule, ctx echo.Context) string {
	if rule == nil {
		return ""
	}
	var saveFile string = rule.SaveFile
	if rule.SaveFunc != nil {
		saveFile = rule.SaveFunc(saveFile, ctx)
	}
	return saveFile
}

func (c *Config) Expired(rule *Rule, ctx echo.Context, saveFile string) (int64, bool) {
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
		case func(string, echo.Context) (int64, bool):
			fn := c.HtmlCacheTime.(func(string, echo.Context) (int64, bool))
			return fn(saveFile, ctx)
		}
	}
	mtime, err := com.FileMTime(saveFile)
	if err != nil {
		ctx.X().Echo().Logger().Debug(err)
	}
	if mtime == 0 {
		return mtime, true
	}
	if time.Now().Local().Unix() > mtime+expired {
		return mtime, true
	}
	return mtime, false
}

func (c *Config) Middleware() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if c.Read(ctx) {
				return nil
			}
			ctx.Set(`webx:ignoreRender`, true)
			if err := h(ctx); err != nil {
				return err
			}
			format, _ := ctx.Get(`webx:format`).(string)
			switch format {
			case `xml`:
				b, err := RenderXML(ctx)
				if err != nil {
					return err
				}
				c.Write(b, ctx)
				return OutputXML(b, ctx)
			case `json`:
				b, err := RenderJSON(ctx)
				if err != nil {
					return err
				}
				c.Write(b, ctx)
				return OutputJSON(b, ctx)
			default:
				b, err := RenderHTML(ctx)
				if err != nil {
					return err
				}
				c.Write(b, ctx)
				return OutputHTML(b, ctx)
			}
		}
	}
}
