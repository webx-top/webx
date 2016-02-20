/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package language

import (
	"regexp"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	X "github.com/webx-top/webx"
	"github.com/webx-top/webx/lib/i18n"
)

const LANG_KEY = `webx:language`

func NewLanguage() *Language {
	return &Language{
		List:     make(map[string]bool),
		Index:    make([]string, 0),
		Default:  "zh-cn",
		uaRegexp: regexp.MustCompile(`;q=[0-9.]+`),
	}
}

type Language struct {
	List     map[string]bool //语种列表
	Index    []string        //索引
	Default  string          //默认语种
	uaRegexp *regexp.Regexp
}

func (a *Language) Set(lang string, on bool, args ...bool) *Language {
	if a.List == nil {
		a.List = make(map[string]bool)
	}
	if _, ok := a.List[lang]; !ok {
		a.Index = append(a.Index, lang)
	}
	a.List[lang] = on
	if on && len(args) > 0 && args[0] {
		a.Default = lang
	}
	return a
}

func (a *Language) IsOk(lang string) bool {
	if on, ok := a.List[lang]; !ok {
		return false
	} else {
		return on
	}
}

func (a *Language) DetectURI(_ engine.Response, r engine.Request) string {
	p := strings.TrimPrefix(r.URL().Path(), `/`)
	s := strings.Index(p, `/`)
	var lang string
	if s != -1 {
		lang = p[0:s]
	} else {
		lang = p
	}
	if lang != "" {
		if on, ok := a.List[lang]; ok {
			r.URL().SetPath(strings.TrimPrefix(p, lang))
			if !on {
				lang = ""
			}
		} else {
			lang = ""
		}
	}
	if lang == "" {
		lang = a.DetectUA(r)
	}
	return lang
}

func (a *Language) DetectUA(r engine.Request) string {
	ua := r.UserAgent()
	ua = a.uaRegexp.ReplaceAllString(ua, ``)
	lg := strings.SplitN(ua, `,`, 5)
	for _, lang := range lg {
		lang = strings.ToLower(lang)
		if a.IsOk(lang) {
			return lang
		}
	}
	return a.Default
}

func (a *Language) Middleware() echo.MiddlewareFunc {
	return echo.MiddlewareFunc(func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			lang := a.DetectURI(c.Response(), c.Request())
			c.SetFunc("Lang", func() string {
				return lang
			})
			c.SetFunc("T", func(key string, args ...interface{}) string {
				return i18n.T(lang, key, args...)
			})
			X.X(c).Language = lang
			return h.Handle(c)
		})
	})
}
