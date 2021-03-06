package pongo2

import (
	"bytes"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"

	. "github.com/admpub/pongo2"
	"github.com/labstack/gommon/log"
	"github.com/webx-top/webx/lib/tplex"
)

func init() {
	tplex.Reg(`pongo2`, func(tmplDir string) tplex.TemplateEx {
		return New(tmplDir)
	})
}

func New(templateDir string) tplex.TemplateEx {
	a := &templatePongo2{
		templateDir: templateDir,
		ext:         `.html`,
		Logger:      log.New("tplex"),
	}
	a.templateDir, _ = filepath.Abs(templateDir)
	return a
}

type templatePongo2 struct {
	templates   map[string]*Template
	mutex       *sync.RWMutex
	loader      TemplateLoader
	set         *TemplateSet
	ext         string
	templateDir string
	Mgr         *tplex.TemplateMgr
	Logger      *log.Logger
	getFuncs    func() map[string]interface{}

	onChange func(string)
}

type templateLoader struct {
	templateDir string
	mgr         *tplex.TemplateMgr
	ext         string
	logger      *log.Logger
}

func (a *templateLoader) Abs(base, name string) string {
	//a.logger.Info(base+" => %v\n", name)
	return filepath.Join(``, name)
}

// Get returns an io.Reader where the template's content can be read from.
func (a *templateLoader) Get(tmpl string) (io.Reader, error) {
	var b []byte
	var e error
	tmpl += a.ext
	if a.mgr != nil && a.mgr.Caches != nil {
		k := strings.TrimPrefix(tmpl, a.templateDir)
		b, e = a.mgr.GetTemplate(k)
	}
	if b == nil || e != nil {
		if e != nil {
			a.logger.Error(e)
		}
		b, e = ioutil.ReadFile(tmpl)
	}
	buf := new(bytes.Buffer)
	buf.WriteString(string(b))
	return buf, e
}

func (a *templatePongo2) MonitorEvent(fn func(string)) {
	a.onChange = fn
}

func (a *templatePongo2) Init(cached ...bool) {
	a.Logger.SetLevel(log.INFO)
	a.Mgr = new(tplex.TemplateMgr)
	a.templates = map[string]*Template{}
	a.mutex = &sync.RWMutex{}
	loader := &templateLoader{
		templateDir: a.templateDir,
		mgr:         a.Mgr,
		ext:         a.ext,
		logger:      a.Logger,
	}
	a.loader = loader
	a.set = NewSet(a.templateDir, a.loader)

	ln := len(cached)
	if ln < 1 || !cached[0] {
		return
	}
	reloadTemplates := true
	if ln > 1 {
		reloadTemplates = cached[1]
	}

	a.Mgr.OnChangeCallback = a.OnChange
	a.Mgr.Init(a.Logger, a.templateDir, reloadTemplates, "*"+a.ext)
}

func (a *templatePongo2) OnChange(name, typ, event string) {
	switch event {
	case "create":
	case "delete", "modify", "rename":
		if typ == "dir" || !strings.HasSuffix(name, a.ext) {
			return
		}
		key := strings.TrimSuffix(name, a.ext)
		//布局模板被修改时，清空缓存
		if strings.HasSuffix(key, `layout`) {
			a.templates = make(map[string]*Template)
			a.Logger.Info(`remove all cached template object:`, name)
		} else if _, ok := a.templates[key]; ok {
			delete(a.templates, key)
			a.Logger.Info(`remove cached template object:`, name)
		}
		if a.onChange != nil {
			a.onChange(name)
		}
	}
}

func (a *templatePongo2) SetFuncMapFn(fn func() map[string]interface{}) {
	a.getFuncs = fn
}

func (a *templatePongo2) Render(w io.Writer, tmpl string, data interface{}, funcMap map[string]interface{}) error {
	t, context := a.parse(tmpl, data, funcMap)
	return t.ExecuteWriter(context, w)
}

func (a *templatePongo2) parse(tmpl string, data interface{}, funcMap map[string]interface{}) (*Template, Context) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	k := tmpl
	if tmpl[0] == '/' {
		k = tmpl[1:]
	}
	t, ok := a.templates[k]
	if !ok {
		var err error
		t, err = a.set.FromFile(tmpl)
		if err != nil {
			t = Must(a.set.FromString(err.Error()))
			return t, Context{}
		}
		a.templates[k] = t
	}
	var context Context
	if a.getFuncs != nil {
		context = Context(a.getFuncs())
	}
	if v, ok := data.(Context); ok {
		if context == nil {
			context = v
		} else {
			for n, f := range v {
				context[n] = f
			}
		}
	} else if v, ok := data.(map[string]interface{}); ok {
		if context == nil {
			context = v
		} else {
			for n, f := range v {
				context[n] = f
			}
		}
	} else {
		if context == nil {
			context = Context{
				`value`: data,
			}
		} else {
			context[`value`] = data
		}
	}
	if funcMap != nil {
		for name, function := range funcMap {
			context[name] = function
			//a.Logger.Info("added func: %v => %#v", name, function)
		}
	}
	return t, context
}

func (a *templatePongo2) Fetch(tmpl string, data interface{}, funcMap map[string]interface{}) string {
	t, context := a.parse(tmpl, data, funcMap)
	r, err := t.Execute(context)
	if err != nil {
		r = err.Error()
	}
	return r
}

func (a *templatePongo2) RawContent(tmpl string) (b []byte, e error) {
	if a.Mgr != nil && a.Mgr.Caches != nil {
		b, e = a.Mgr.GetTemplate(tmpl)
	}
	if b == nil || e != nil {
		b, e = ioutil.ReadFile(filepath.Join(a.templateDir, tmpl))
	}
	return
}

func (a *templatePongo2) ClearCache() {
	if a.Mgr != nil {
		a.Mgr.ClearCache()
	}
	a.templates = make(map[string]*Template)
}

func (a *templatePongo2) Close() {
	a.ClearCache()
	if a.Mgr != nil {
		a.Mgr.Close()
	}
}
