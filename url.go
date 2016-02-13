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
package webx

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"

	"github.com/webx-top/webx/lib/com"
)

func NewURL(project string, serv *Server) *URL {
	return &URL{
		projectPath: `github.com/webx-top/` + project,
		urls:        make(map[string]*Url),
		Server:      serv,
	}
}

type URL struct {
	projectPath string
	urls        map[string]*Url
	*Server
}

func (a *URL) SetProjectPath(projectPath string) {
	a.projectPath = strings.TrimSuffix(projectPath, `/`)
}

func (a *URL) Build(app string, ctl string, act string, params ...interface{}) (r string) {
	pkg := a.projectPath + `/app/` + app + `/controller`
	key := ``
	if ctl == `` {
		key = pkg + `.` + act
	} else {
		key = pkg + `.(*` + ctl + `).` + act + `-fm`
	}
	if u, ok := a.urls[key]; ok {
		r = u.Gen(params)
	}
	r = strings.TrimSuffix(a.Server.App(app).Url, `/`) + r
	return
}

func (a *URL) BuildByPath(path string, args ...map[string]interface{}) (r string) {
	var app, ctl, act string
	uris := strings.SplitN(path, "?", 2)
	ret := strings.SplitN(uris[0], `/`, 3)
	switch len(ret) {
	case 3:
		act = ret[2]
		ctl = ret[1]
		app = ret[0]
	case 2:
		act = ret[1]
		app = ret[0]
	default:
		return
	}
	pkg := a.projectPath + `/app/` + app + `/controller`
	key := ``
	if ctl == `` {
		key = pkg + `.` + act
	} else {
		key = pkg + `.(*` + ctl + `).` + act + `-fm`
	}
	var params url.Values
	if len(uris) > 1 {
		params, _ = url.ParseQuery(uris[1])
	}
	if len(args) > 0 {
		for k, v := range args[0] {
			params.Set(k, fmt.Sprintf("%v", v))
		}
	}
	if u, ok := a.urls[key]; ok {
		r = u.Gen(params)
	}
	r = strings.TrimSuffix(a.Server.App(app).Url, `/`) + r
	return
}

func (a *URL) FuncPath(h interface{}) string {
	return com.FuncName(h)
}

func (a *URL) Set(route string, h interface{}, memo ...string) (pkg string, ctl string, act string) {
	key := a.FuncPath(h)
	a.Server.Core.Logger().Infof(`URL:%v => %v`, route, key)
	urls := &Url{}
	urls.Set(route, memo...)
	a.urls[key] = urls
	pkg, ctl, act = com.ParseFuncName(key)
	return
}

func (a *URL) SetByKey(route string, key string, memo ...string) *Url {
	a.Server.Core.Logger().Infof(`URL:%v => %v`, route, key)
	urls := NewUrl()
	urls.Set(route, memo...)
	a.urls[key] = urls
	return urls
}

func (a *URL) Urls() map[string]*Url {
	return a.urls
}

func (a *URL) Get(key string) (u *Url) {
	u, _ = a.urls[key]
	return
}

type Url struct {
	Route  string
	Format string
	Params []string
	Memo   string
	exts   map[string]int
}

func NewUrl() *Url {
	return &Url{
		Params: []string{},
		exts:   map[string]int{},
	}
}

func (m *Url) SetExts(exts []string) {
	for key, val := range exts {
		m.exts[val] = key
	}
}

func (m *Url) ValidExt(ext string) (ok bool) {
	if len(m.exts) < 1 {
		ok = true
	} else {
		_, ok = m.exts[ext]
	}
	return
}

func (m *Url) Gen(vals interface{}) (r string) {
	r = m.Route
	if r == "" {
		return
	}
	switch vals.(type) {
	case url.Values:
		val := vals.(url.Values)
		for _, name := range m.Params {
			tag := `:` + name
			v := val.Get(name)
			r = strings.Replace(r, tag+`/`, v+`/`, -1)
			if strings.HasSuffix(r, tag) {
				r = strings.TrimSuffix(r, tag) + v
			}
			val.Del(name)
		}
		q := val.Encode()
		if q != `` {
			r += `?` + q
		}
	case map[string]string:
		val := vals.(map[string]string)
		for _, name := range m.Params {
			tag := `:` + name
			v, _ := val[name]
			r = strings.Replace(r, tag+`/`, v+`/`, -1)
			if strings.HasSuffix(r, tag) {
				r = strings.TrimSuffix(r, tag) + v
			}
		}
	case []interface{}:
		val := vals.([]interface{})
		r = fmt.Sprintf(m.Format, val...)
	default:
	}
	return
}

func (m *Url) Set(route string, memo ...string) {
	m.Route = route
	m.Params = make([]string, 0)
	uri := new(bytes.Buffer)
	for i, l := 0, len(route); i < l; i++ {
		if route[i] == ':' {
			start := i + 1
			for ; i < l && route[i] != '/'; i++ {
			}
			if i > start {
				m.Params = append(m.Params, route[start:i])
			}
			uri.WriteString("%v")
		}
		if i < l {
			uri.WriteByte(route[i])
		}
	}
	m.Format = uri.String()
	if len(memo) > 0 {
		m.Memo = memo[0]
	}
}
