package middleware

import (
	"regexp"
	"strings"

	"github.com/admpub/echo"
)

func NewLanguage() *Language {
	return &Language{
		List:     make(map[string]bool),
		Index:    make([]string, 0),
		Default:  "",
		uaRegexp: regexp.MustCompile(`;q=[0-9.]+`),
	}
}

type Language struct {
	List     map[string]bool //语种列表
	Index    []string        //索引
	Default  string          //默认语种
	uaRegexp *regexp.Regexp
}

func (a *Language) Set(lang string, on bool) *Language {
	if a.List == nil {
		a.List = make(map[string]bool)
	}
	if _, ok := a.List[lang]; !ok {
		a.Index = append(a.Index, lang)
	}
	a.List[lang] = on
	return a
}

func (a *Language) IsOk(lang string) bool {
	if on, ok := a.List[lang]; !ok {
		return false
	} else {
		return on
	}
}

func (a *Language) DetectURI() echo.HandlerFunc {
	return func(c *echo.Context) error {
		p := strings.TrimPrefix(c.Request().URL.Path, `/`)
		s := strings.Index(p, `/`)
		var lang string
		if s > 0 {
			lang = p[0:s]
			if on, ok := a.List[lang]; ok {
				c.Request().URL.Path = p[s:]
				if on {
					c.Set("language", lang)
				} else {
					lang = ""
				}
			} else {
				lang = ""
			}
		}
		if lang == "" {
			a.DetectUA(c)
		}
		return nil
	}
}

func (a *Language) DetectUA(c *echo.Context) *Language {
	ua := c.Request().UserAgent()
	ua = a.uaRegexp.ReplaceAllString(ua, ``)
	lg := strings.SplitN(ua, `,`, 5)
	for _, lang := range lg {
		lang = strings.ToLower(lang)
		if a.IsOk(lang) {
			c.Set("language", lang)
			return a
		}
	}
	c.Set("language", a.Default)
	return a
}
