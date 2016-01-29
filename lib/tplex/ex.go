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

/**
 * 模板扩展
 * @author swh <swh@admpub.com>
 */
package tplex

import (
	"bytes"
	"errors"
	"fmt"
	htmlTpl "html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/labstack/gommon/log"
)

var Debug = false

func New(templateDir string) *TemplateEx {
	t := &TemplateEx{
		CachedRelation: make(map[string]*CcRel),
		TemplateDir:    templateDir,
		DelimLeft:      "{{",
		DelimRight:     "}}",
		IncludeTag:     "Include",
		ExtendTag:      "Extend",
		BlockTag:       "Block",
		SuperTag:       "Super",
		Ext:            ".html",
		Debug:          Debug,
		Cached:         true,
	}
	t.Logger = log.New("tplex")
	t.Logger.SetLevel(log.INFO)
	t.InitRegexp()
	return t
}

type CcRel struct {
	Rel  map[string]uint8
	Tpl  [2]*htmlTpl.Template //0是独立模板；1是子模板
	Sub  string
	Self string
}

type TemplateEx struct {
	CachedRelation     map[string]*CcRel
	TemplateDir        string
	TemplateMgr        *TemplateMgr
	BeforeRender       func(*string)
	DelimLeft          string
	DelimRight         string
	incTagRegex        *regexp.Regexp
	extTagRegex        *regexp.Regexp
	blkTagRegex        *regexp.Regexp
	cachedRegexIdent   string
	IncludeTag         string
	ExtendTag          string
	BlockTag           string
	SuperTag           string
	Ext                string
	TemplatePathParser func(string) string
	Debug              bool
	Cached             bool
	FuncMapFn          func() htmlTpl.FuncMap
	Logger             *log.Logger
	FileChangeEvent    func(string)
}

func (self *TemplateEx) InitMgr(cached ...bool) {
	self.TemplateMgr = new(TemplateMgr)

	ln := len(cached)
	if ln < 1 || !cached[0] {
		return
	}
	reloadTemplates := true
	if ln > 1 {
		reloadTemplates = cached[1]
	}

	self.TemplateMgr.OnChangeCallback = func(name, typ, event string) {
		switch event {
		case "create":
		case "delete", "modify", "rename":
			if typ == "dir" {
				return
			}
			if cs, ok := self.CachedRelation[name]; ok {
				for key, _ := range cs.Rel {
					if name == key {
						self.Logger.Info("remove cached template object: %v", key)
						continue
					}
					if _, ok := self.CachedRelation[key]; ok {
						self.Logger.Info("remove cached template object: %v", key)
						delete(self.CachedRelation, key)
					}
				}
				delete(self.CachedRelation, name)
			}
			if self.FileChangeEvent != nil {
				self.FileChangeEvent(name)
			}
		}
	}
	self.TemplateMgr.Init(self.Logger, self.TemplateDir, reloadTemplates, "*"+self.Ext)
}

func (self *TemplateEx) SetMgr(mgr *TemplateMgr) {
	self.TemplateMgr = mgr
}

func (self *TemplateEx) TemplatePath(p string) string {
	if self.TemplatePathParser == nil {
		return p
	}
	return self.TemplatePathParser(p)
}

func (self *TemplateEx) echo(messages ...string) {
	if self.Debug {
		var message string
		for _, v := range messages {
			message += v + ` `
		}
		fmt.Println(`[tplex]`, message)
	}
}

func (self *TemplateEx) InitRegexp() {
	left := regexp.QuoteMeta(self.DelimLeft)
	right := regexp.QuoteMeta(self.DelimRight)
	rfirst := regexp.QuoteMeta(self.DelimRight[0:1])
	self.incTagRegex = regexp.MustCompile(left + self.IncludeTag + `[\s]+"([^"]+)"(?:[\s]+([^` + rfirst + `]+))?[\s]*` + right)
	self.extTagRegex = regexp.MustCompile(left + self.ExtendTag + `[\s]+"([^"]+)"(?:[\s]+([^` + rfirst + `]+))?[\s]*` + right)
	self.blkTagRegex = regexp.MustCompile(`(?s)` + left + self.BlockTag + `[\s]+"([^"]+)"[\s]*` + right + `(.*?)` + left + `\/` + self.BlockTag + right)
}

// Render HTML
func (self *TemplateEx) Render(w io.Writer, tmplName string, values interface{}, funcs htmlTpl.FuncMap) error {
	var funcMap htmlTpl.FuncMap
	if self.FuncMapFn != nil {
		funcMap = self.FuncMapFn()
		if funcs != nil {
			for k, v := range funcs {
				funcMap[k] = v
			}
		}
	} else {
		if funcs != nil {
			funcMap = funcs
		}
	}
	tmpl := self.Fetch(tmplName, funcMap)
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, tmpl.Name(), values)
	if err != nil {
		return errors.New(fmt.Sprintf("Parse %v err: %v", tmpl.Name(), err))
	}
	_, err = io.Copy(w, buf)
	if err != nil {
		return errors.New(fmt.Sprintf("Parse %v err: %v", tmpl.Name(), err))
	}
	return err
}

func (self *TemplateEx) Fetch(tmplName string, funcMap htmlTpl.FuncMap) (tmpl *htmlTpl.Template) {
	tmplName = tmplName + self.Ext
	tmplName = self.TemplatePath(tmplName)
	rel, ok := self.CachedRelation[tmplName]
	if !ok || rel.Tpl[0] == nil {
		if rel == nil {
			rel = &CcRel{
				Rel: map[string]uint8{tmplName: 0},
				Tpl: [2]*htmlTpl.Template{},
			}
		}
		self.echo(`Read not cached template content:`, tmplName)
		b, err := self.RawContent(tmplName)
		if err != nil {
			self.Logger.Error("RenderTemplate %v read err: %s", tmplName, err)
		}

		content := string(b)
		if self.BeforeRender != nil {
			self.BeforeRender(&content)
		}
		subcs := make(map[string]string, 0) //子模板内容
		extcs := make(map[string]string, 0) //母板内容

		ident := self.DelimLeft + self.IncludeTag + self.DelimRight
		if self.cachedRegexIdent != ident || self.incTagRegex == nil {
			self.InitRegexp()
		}
		m := self.extTagRegex.FindAllStringSubmatch(content, 1)
		if len(m) > 0 {
			self.ParseBlock(content, &subcs, &extcs)
			extFile := m[0][1] + self.Ext
			passObject := m[0][2]
			extFile = self.TemplatePath(extFile)
			self.echo(`Read layout template content:`, extFile)
			b, err = self.RawContent(extFile)
			if err != nil {
				content = fmt.Sprintf("RenderTemplate %v read err: %s", extFile, err)
			} else {
				content = string(b)
			}
			content = self.ParseExtend(content, &extcs, passObject, &subcs)

			if v, ok := self.CachedRelation[extFile]; !ok {
				self.CachedRelation[extFile] = &CcRel{
					Rel: map[string]uint8{tmplName: 0},
					Tpl: [2]*htmlTpl.Template{},
				}
			} else if _, ok := v.Rel[tmplName]; !ok {
				self.CachedRelation[extFile].Rel[tmplName] = 0
			}
		}
		content = self.ContainsSubTpl(content, &subcs)
		t := htmlTpl.New(tmplName)
		t.Delims(self.DelimLeft, self.DelimRight)
		t.Funcs(funcMap)
		self.echo(`The template content:`, content)
		tmpl, err = t.Parse(content)
		if err != nil {
			content = fmt.Sprintf("Parse %v err: %v", tmplName, err)
			tmpl, _ = t.Parse(content)
		}
		for name, subc := range subcs {
			v, ok := self.CachedRelation[name]
			if ok && v.Tpl[1] != nil {
				self.CachedRelation[name].Rel[tmplName] = 0
				tmpl.AddParseTree(name, self.CachedRelation[name].Tpl[1].Tree)
				continue
			}
			if self.BeforeRender != nil {
				self.BeforeRender(&subc)
			}
			var t *htmlTpl.Template
			if name == tmpl.Name() {
				t = tmpl
			} else {
				t = tmpl.New(name)
				_, err = t.Parse(subc)
				if err != nil {
					t.Parse(fmt.Sprintf("Parse File %v err: %v", name, err))
				}
			}

			if ok {
				self.CachedRelation[name].Rel[tmplName] = 0
				self.CachedRelation[name].Tpl[1] = t
			} else {
				self.CachedRelation[name] = &CcRel{
					Rel: map[string]uint8{tmplName: 0},
					Tpl: [2]*htmlTpl.Template{nil, t},
				}
			}

		}
		for name, extc := range extcs {
			if self.BeforeRender != nil {
				self.BeforeRender(&extc)
			}
			var t *htmlTpl.Template
			if name == tmpl.Name() {
				t = tmpl
			} else {
				t = tmpl.New(name)
				_, err = t.Parse(extc)
				if err != nil {
					t.Parse(fmt.Sprintf("Parse Block %v err: %v", name, err))
				}
			}
		}

		rel.Tpl[0] = tmpl
		self.CachedRelation[tmplName] = rel

	} else {
		tmpl = rel.Tpl[0]
		tmpl.Funcs(funcMap)
		if self.Debug {
			fmt.Println(`Using the template object to be cached:`, tmplName)
			fmt.Println("_________________________________________")
			fmt.Println("")
			for k, v := range tmpl.Templates() {
				fmt.Printf("%v. %#v\n", k, v.Name())
			}
			fmt.Println("_________________________________________")
			fmt.Println("")
		}
	}
	return tmpl
}

func (self *TemplateEx) ParseBlock(content string, subcs *map[string]string, extcs *map[string]string) {
	matches := self.blkTagRegex.FindAllStringSubmatch(content, -1)
	for _, v := range matches {
		blockName := v[1]
		content := v[2]
		(*extcs)[blockName] = self.Tag(`define "`+blockName+`"`) + self.ContainsSubTpl(content, subcs) + self.Tag(`end`)
	}
}

func (self *TemplateEx) ParseExtend(content string, extcs *map[string]string, passObject string, subcs *map[string]string) string {
	if passObject == "" {
		passObject = "."
	}
	matches := self.blkTagRegex.FindAllStringSubmatch(content, -1)
	var superTag string
	if self.SuperTag != "" {
		superTag = self.Tag(self.SuperTag)
	}
	var rec map[string]uint8 = make(map[string]uint8)
	for _, v := range matches {
		matched := v[0]
		blockName := v[1]
		innerStr := v[2]
		if v, ok := (*extcs)[blockName]; ok {
			rec[blockName] = 0
			if superTag != "" && strings.Contains(v, superTag) {
				innerStr = self.ContainsSubTpl(innerStr, subcs)
				(*extcs)[blockName] = strings.Replace(v, superTag, innerStr, 1)
			}
			content = strings.Replace(content, matched, self.Tag(`template "`+blockName+`" `+passObject), -1)
		} else {
			content = strings.Replace(content, matched, innerStr, -1)
		}
	}
	for k, _ := range *extcs {
		if _, ok := rec[k]; !ok {
			delete(*extcs, k)
		}
	}
	return content
}

func (self *TemplateEx) ContainsSubTpl(content string, subcs *map[string]string) string {
	matches := self.incTagRegex.FindAllStringSubmatch(content, -1)
	for _, v := range matches {
		matched := v[0]
		tmplFile := v[1]
		passObject := v[2]
		tmplFile += self.Ext
		tmplFile = self.TemplatePath(tmplFile)
		if _, ok := (*subcs)[tmplFile]; !ok {
			if v, ok := self.CachedRelation[tmplFile]; ok && v.Tpl[1] != nil {
				(*subcs)[tmplFile] = ""
			} else {
				b, err := self.RawContent(tmplFile)
				if err != nil {
					return fmt.Sprintf("RenderTemplate %v read err: %s", tmplFile, err)
				}
				str := string(b)
				(*subcs)[tmplFile] = "" //先登记，避免死循环
				str = self.ContainsSubTpl(str, subcs)
				(*subcs)[tmplFile] = self.Tag(`define "`+tmplFile+`"`) + str + self.Tag(`end`)
			}
		}
		if passObject == "" {
			passObject = "."
		}
		content = strings.Replace(content, matched, self.Tag(`template "`+tmplFile+`" `+passObject), -1)
	}
	return content
}

func (self *TemplateEx) Tag(content string) string {
	return self.DelimLeft + content + self.DelimRight
}

func (self *TemplateEx) Include(tmplName string, funcMap htmlTpl.FuncMap, values interface{}) interface{} {
	tmpl := self.Fetch(tmplName, funcMap)
	return htmlTpl.HTML(self.Parse(tmpl, values))
}

func (self *TemplateEx) Parse(tmpl *htmlTpl.Template, values interface{}) string {
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, tmpl.Name(), values)
	if err != nil {
		return fmt.Sprintf("Parse %v err: %v", tmpl.Name(), err)
	}
	b, err := ioutil.ReadAll(buf)
	if err != nil {
		return fmt.Sprintf("Parse %v err: %v", tmpl.Name(), err)
	}
	return string(b)
}

func (self *TemplateEx) RawContent(tmpl string) ([]byte, error) {
	if self.TemplateMgr != nil && self.TemplateMgr.Caches != nil {
		return self.TemplateMgr.GetTemplate(tmpl)
	}
	return ioutil.ReadFile(filepath.Join(self.TemplateDir, tmpl))
}

func (self *TemplateEx) ClearCache() {
	if self.TemplateMgr != nil {
		self.TemplateMgr.ClearCache()
	}
	self.CachedRelation = make(map[string]*CcRel)
}

func (self *TemplateEx) Close() {
	self.ClearCache()
	if self.TemplateMgr != nil {
		self.TemplateMgr.Close()
	}
}
