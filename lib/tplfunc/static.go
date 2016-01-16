package tplfunc

import (
	"fmt"
	"html/template"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/webx-top/webx/lib/com"
	"github.com/webx-top/webx/lib/minify"
)

var (
	regexCssUrlAttr      *regexp.Regexp = regexp.MustCompile(`url\(['"]?(\.\./[^\)'"]+)['"]?\)`)
	regexCssImport       *regexp.Regexp = regexp.MustCompile(`@import[\s]+["']([^"']+)["']`)
	regexCssCleanSpace   *regexp.Regexp = regexp.MustCompile(`(?s)\s*(\{|\}|;|:)\s*`)
	regexCssCleanSpace2  *regexp.Regexp = regexp.MustCompile(`(?s)\s{2,}`)
	regexCssCleanComment *regexp.Regexp = regexp.MustCompile(`(?s)[\s]*/\*(.*?)\*/[\s]*`)
)

func NewStatic(staticPath, rootPath string) *Static {
	return &Static{
		Path:            staticPath,
		RootPath:        rootPath,
		CombineJs:       true,
		CombineCss:      true,
		CombineSavePath: `combine`,
		Combined:        make(map[string][]string),
		Combines:        make(map[string]bool),
		mutex:           &sync.Mutex{},
	}
}

type Static struct {
	RootPath        string //根路径：相对于本程序的路径,本程序读取时需要
	Path            string //网址访问的路径
	CombineJs       bool
	CombineCss      bool
	CombineSavePath string //合并文件保存路径，首尾均不带斜杠
	Combined        map[string][]string
	Combines        map[string]bool
	mutex           *sync.Mutex
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
	if len(staticFiles) == 1 || !s.CombineJs {
		for _, staticFile := range staticFiles {
			r += `<script type="text/javascript" src="` + s.JsUrl(staticFile) + `" charset="utf-8"></script>`
		}
		return template.HTML(r)
	}

	r = s.CombineSavePath + "/" + com.Md5(strings.Join(staticFiles, "|")) + ".js"
	if s.IsCombined(r) == false || com.FileExists(s.RootPath+"/"+r) == false {
		var content string
		for _, url := range staticFiles {
			urlFile := s.RootPath + "/js/" + url
			if con, err := com.ReadFileS(urlFile); err != nil {
				fmt.Println(err)
			} else {
				s.RecordCombined("js/"+url, r)
				content += "\n/* <from: " + url + "> */\n"
				b, err := minify.MinifyJS([]byte(con))
				if err != nil {
					fmt.Println(err)
				}
				con = string(b)
				con = regexCssCleanComment.ReplaceAllString(con, ``)
				content += con
			}
			//fmt.Println(url)
		}
		com.WriteFile(s.RootPath+"/"+r, []byte(content))
		s.RecordCombines(r)
	}
	r = `<script type="text/javascript" src="` + s.StaticUrl(r) + `" charset="utf-8"></script>`
	return template.HTML(r)
}

func (s *Static) CssTag(staticFiles ...string) template.HTML {
	var r string
	if len(staticFiles) == 1 || !s.CombineCss {
		for _, staticFile := range staticFiles {
			r += `<link rel="stylesheet" type="text/css" href="` + s.CssUrl(staticFile) + `" charset="utf-8" />`
		}
		return template.HTML(r)
	}

	r = s.CombineSavePath + "/" + com.Md5(strings.Join(staticFiles, "|")) + ".css"
	if s.IsCombined(r) == false || com.FileExists(s.RootPath+"/"+r) == false {
		var content string
		for _, url := range staticFiles {
			urlFile := s.RootPath + "/css/" + url
			if con, err := com.ReadFileS(urlFile); err != nil {
				fmt.Println(err)
			} else {
				all := regexCssUrlAttr.FindAllStringSubmatch(con, -1)
				dir := path.Dir(s.CssUrl(url))
				for _, v := range all {
					res := dir
					val := v[1]
					for strings.HasPrefix(val, "../") {
						res = path.Dir(res)
						val = strings.TrimPrefix(val, "../")
					}
					con = strings.Replace(con, v[0], "url('"+res+"/"+strings.TrimLeft(val, "/")+"')", 1)
				}
				all = regexCssImport.FindAllStringSubmatch(con, -1)
				for _, v := range all {
					res := dir
					val := v[1]
					for strings.HasPrefix(val, "../") {
						res = path.Dir(res)
						val = strings.TrimPrefix(val, "../")
					}
					con = strings.Replace(con, v[0], `@import "`+res+"/"+strings.TrimLeft(val, "/")+`"`, 1)
				}
				s.RecordCombined("css/"+url, r)
				content += "\n/* <from: " + url + "> */\n"
				con = regexCssCleanComment.ReplaceAllString(con, ``)
				con = regexCssCleanSpace.ReplaceAllString(con, `$1`)
				con = regexCssCleanSpace2.ReplaceAllString(con, ` `)
				content += con
			}
		}
		com.WriteFile(s.RootPath+"/"+r, []byte(content))
		s.RecordCombines(r)
	}
	r = `<link rel="stylesheet" type="text/css" href="` + s.StaticUrl(r) + `" charset="utf-8" />`
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

func (s *Static) DeleteCombined(url string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if val, ok := s.Combined[url]; ok {
		for _, v := range val {
			if _, has := s.Combines[v]; !has {
				continue
			}
			err := os.Remove(filepath.Join(s.RootPath, v))
			delete(s.Combines, v)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (s *Static) RecordCombined(fromUrl string, combineUrl string) {
	if s.Combined == nil {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if _, ok := s.Combined[fromUrl]; !ok {
		s.Combined[fromUrl] = make([]string, 0)
	}
	s.Combined[fromUrl] = append(s.Combined[fromUrl], combineUrl)
}

func (s *Static) RecordCombines(combineUrl string) {
	if s.Combines == nil {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Combines[combineUrl] = true
}

func (s *Static) IsCombined(combineUrl string) (ok bool) {
	if s.Combines == nil {
		return
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, ok = s.Combines[combineUrl]
	return
}

func (s *Static) ClearCache() {
	for f, _ := range s.Combines {
		os.Remove(filepath.Join(s.RootPath, f))
	}
	s.Combined = make(map[string][]string)
	s.Combines = make(map[string]bool)
}
