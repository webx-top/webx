package language

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/context"
	"github.com/webx-top/echo"
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

func (a *Language) DetectURI(_ http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, `/`)
	s := strings.Index(p, `/`)
	var lang string
	if s != -1 {
		lang = p[0:s]
	} else {
		lang = p
	}
	if lang != "" {
		if on, ok := a.List[lang]; ok {
			r.URL.Path = strings.TrimPrefix(p, lang)
			if on {
				context.Set(r, LANG_KEY, lang)
			} else {
				lang = ""
			}
		} else {
			lang = ""
		}
	}
	if lang == "" {
		a.DetectUA(r)
	}
}

func (a *Language) DetectUA(r *http.Request) *Language {
	ua := r.UserAgent()
	ua = a.uaRegexp.ReplaceAllString(ua, ``)
	lg := strings.SplitN(ua, `,`, 5)
	for _, lang := range lg {
		lang = strings.ToLower(lang)
		if a.IsOk(lang) {
			context.Set(r, LANG_KEY, lang)
			return a
		}
	}
	context.Set(r, LANG_KEY, a.Default)
	return a
}

//存储到echo.Context中
func (a *Language) Store() echo.HandlerFunc {
	return func(c echo.Context) error {
		lang := context.Get(c.Request(), LANG_KEY)
		language, _ := lang.(string)
		//TODO: 移到echo.Context中
		if language == `` {
			language = a.Default
		}
		c.SetFunc("Lang", func() string {
			return language
		})
		c.SetFunc("T", func(key string, args ...interface{}) string {
			return i18n.T(language, key, args...)
		})
		X.X(c).Language = language
		return nil
	}
}
