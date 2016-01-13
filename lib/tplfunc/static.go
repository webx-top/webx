package tplfunc

import (
	"html/template"
)

func NewStatic(staticPath string) *Static {
	return &Static{Path: staticPath}
}

type Static struct {
	Path string
}

func (s *Static) StaticUrl(staticFile string) (r string) {
	r = s.Path + "/" + staticFile
	return
}

func (s *Static) JsUrl(staticFile string) (r string) {
	r = s.StaticUrl("js/" + staticFile)
	return
}

func (s *Static) CssUrl(staticFile string) (r string) {
	r = s.StaticUrl("css/" + staticFile)
	return r
}

func (s *Static) ImgUrl(staticFile string) (r string) {
	r = s.StaticUrl("img/" + staticFile)
	return r
}

func (s *Static) JsTag(staticFiles ...string) template.HTML {
	var r string
	for _, staticFile := range staticFiles {
		r += `<script type="text/javascript" src="` + s.JsUrl(staticFile) + `"></script>`
	}
	return template.HTML(r)
}

func (s *Static) CssTag(staticFiles ...string) template.HTML {
	var r string
	for _, staticFile := range staticFiles {
		r += `<link rel="stylesheet" href="` + s.CssUrl(staticFile) + `" />`
	}
	return template.HTML(r)
}

func (s *Static) ImgTag(staticFile string, attrs ...string) template.HTML {
	var attr string
	for i, l := 0, len(attrs); i+1 < l; i++ {
		var k, v string
		k = attrs[i]
		i++
		v = attrs[i]
		attr += ` ` + k + `="` + v + `"`
	}
	r := `<img src="` + s.ImgUrl(staticFile) + `"` + attr + ` />`
	return template.HTML(r)
}

func (s *Static) Register(funcMap template.FuncMap) template.FuncMap {
	funcMap["StaticUrl"] = s.StaticUrl
	funcMap["JsUrl"] = s.JsUrl
	funcMap["CssUrl"] = s.CssUrl
	funcMap["ImgUrl"] = s.ImgUrl
	funcMap["JsTag"] = s.JsTag
	funcMap["CssTag"] = s.CssTag
	funcMap["ImgTag"] = s.ImgTag
	return funcMap
}
