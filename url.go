package webx

import (
	"bytes"
	"fmt"
	"net/url"
	"reflect"
	"runtime"
	"strings"
)

func NewURL(project string, serv *Server) *URL {
	return &URL{
		Project: project,
		urls:    make(map[string]*Url),
		Server:  serv,
	}
}

type URL struct {
	Project string
	urls    map[string]*Url
	*Server
}

func (a *URL) Build(app string, ctl string, act string, params ...interface{}) (r string) {
	pkg := `github.com/webx-top/` + a.Project + `/app/` + app + `/controller`
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
	pkg := `github.com/webx-top/` + a.Project + `/app/` + app + `/controller`
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

func (a *URL) Set(route string, h interface{}, memo ...string) {
	key := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	urls := &Url{}
	urls.Set(route, memo...)
	a.urls[key] = urls
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
			r = strings.Replace(r, `:`+name+`/`, val.Get(name)+`/`, -1)
			val.Del(name)
		}
		q := val.Encode()
		if q != `` {
			r += `?` + q
		}
	case map[string]string:
		val := vals.(map[string]string)
		for _, name := range m.Params {
			v, _ := val[name]
			r = strings.Replace(r, `:`+name+`/`, v+`/`, -1)
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
