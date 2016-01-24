package webx

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/webx-top/echo"
)

var (
	mapperType        = reflect.TypeOf(Mapper{})
	methodSuffixRegex = regexp.MustCompile(`(?:_(?:` + strings.Join(echo.Methods(), `|`) + `))+$`)
)

type Mapper struct{}

type BeforeHandler interface {
	Before(*Context) error
}

type AfterHandler interface {
	After(*Context) error
}

type Initer interface {
	Init(*Context, *App)
}

type Before interface {
	Before() error
}

type After interface {
	After() error
}

type HandlerFunc func(*Context) error

type Controller struct {
	BeforeHandler HandlerFunc
	AfterHandler  HandlerFunc

	HasBefore bool
	HasAfter  bool
	Init      func(*Context, *App)

	Controller interface{}
	Webx       Webxer
	*App
}

func (a *Controller) wrapHandler(h HandlerFunc, ctl string, act string) func(echo.Context) error {
	if a.BeforeHandler != nil && a.AfterHandler != nil {
		return func(ctx echo.Context) error {
			c := X(ctx)
			c.Init(ctl, act)
			if err := a.BeforeHandler(c); err != nil {
				return err
			}
			if c.Exit {
				return nil
			}
			if err := h(c); err != nil {
				return err
			}
			if c.Exit {
				return nil
			}
			return a.AfterHandler(c)
		}
	}
	if a.BeforeHandler != nil {
		return func(ctx echo.Context) error {
			c := X(ctx)
			c.Init(ctl, act)
			if err := a.BeforeHandler(c); err != nil {
				return err
			}
			if c.Exit {
				return nil
			}
			return h(c)
		}
	}
	if a.AfterHandler != nil {
		return func(ctx echo.Context) error {
			c := X(ctx)
			c.Init(ctl, act)
			if err := h(c); err != nil {
				return err
			}
			if c.Exit {
				return nil
			}
			return a.AfterHandler(c)
		}
	}
	return func(ctx echo.Context) error {
		c := X(ctx)
		c.Init(ctl, act)
		return h(c)
	}
}

//注册路由：Controller.R(`/index`,Index.Index,"GET","POST")
func (a *Controller) R(path string, h HandlerFunc, methods ...string) *Controller {
	if len(methods) < 1 {
		methods = append(methods, "GET")
	}
	_, ctl, act := a.App.Server.URL.Set(path, h)
	a.Webx.Match(methods, path, a.wrapHandler(h, ctl, act))
	return a
}

func (a *Controller) RouteByTag() {
	if _, valid := a.Controller.(Initer); !valid {
		a.Server.Echo.Logger().Info("%T is no method Init(*Context,*App),skip.", a.Controller)
		return
	}
	t := reflect.TypeOf(a.Controller)
	e := t.Elem()
	v := reflect.ValueOf(a.Controller)
	ctlPath := e.PkgPath() + ".(*" + e.Name() + ")."
	//github.com/webx-top/{Project}/app/{App}/controller.(*Index).
	ctl := strings.ToLower(e.Name())

	for i := 0; i < e.NumField(); i++ {
		f := e.Field(i)
		if f.Type != mapperType {
			continue
		}
		fn := strings.Title(f.Name)
		name := fn
		m := v.MethodByName(fn)
		if !m.IsValid() {
			continue
		}

		//支持的tag:
		// 1. webx - 路由规则
		// 2. memo - 注释说明
		//webx标签内容支持以下格式：
		// 1、只指定http请求方式，如`webx:"POST|GET"`
		// 2、只指定路由规则，如`webx:"index"`
		// 3、只指定扩展名规则，如`webx:".JSON|XML"`
		// 4、指定以上全部规则，如`webx:"GET|POST.JSON|XML index"`
		tag := e.Field(i).Tag
		tagv := tag.Get("webx")
		methods := []string{}
		extends := []string{}
		var p, w string
		if tagv != "" {
			tags := strings.Split(tagv, " ")
			length := len(tags)
			if length >= 2 { //`webx:"GET|POST /index"`
				w = tags[0]
				p = tags[1]
			} else if length == 1 {
				if matched, _ := regexp.MatchString(`^[A-Z.]+(\|[A-Z]+)*$`, tags[0]); !matched {
					//非全大写字母时，判断为网址规则
					p = tags[0]
				} else { //`webx:"GET|POST"`
					w = tags[0]
				}
			}
		}
		if p == "" {
			p = "/" + f.Name
		} else if p[0] != '/' {
			p = "/" + p
		}
		path := "/" + ctl + "/" + p
		met := ""
		ext := ""
		if w != "" {
			me := strings.Split(w, ".")
			met = me[0]
			if len(me) > 1 {
				ext = me[1]
			}
		}
		if met != "" {
			methods = strings.Split(met, "|")
		}
		if ext != "" {
			extends = strings.Split(ext, "|")
		}
		k := ctlPath + name + "-fm"
		u := a.App.Server.URL.SetByKey(path, k, tag.Get("memo"))
		u.SetExts(extends)
		h := func(ctx echo.Context) error {
			c := X(ctx)
			c.Init(e.Name(), name)
			if !u.ValidExt(c.Format) {
				return c.HTML(404, `The contents can not be displayed in this format: `+c.Format)
			}
			v := reflect.New(e)
			ac := v.Interface()
			ac.(Initer).Init(c, a.App)
			if a.HasBefore {
				if err := ac.(Before).Before(); err != nil {
					return err
				}
				if c.Exit {
					return nil
				}
			}
			m := v.MethodByName(fn + `_` + c.Method() + strings.ToUpper(c.Format))
			if !m.IsValid() {
				m = v.MethodByName(fn + `_` + c.Method())
				if !m.IsValid() {
					m = v.MethodByName(fn)
				}
			}
			if r, err := a.SafelyCall(m, []reflect.Value{}); err != nil {
				return err
			} else if len(r) > 0 {
				if err, ok := r[0].Interface().(error); ok && err != nil {
					return err
				}
			}
			if a.HasAfter {
				if c.Exit {
					return nil
				}
				return ac.(After).After()
			}
			return nil
		}
		if len(met) < 1 {
			a.Webx.Any(path, h)
			for strings.HasSuffix(path, `/index`) {
				path = strings.TrimSuffix(path, `/index`)
				a.Webx.Any(path+`/`, h)
			}
			continue
		}
		a.Webx.Match(methods, path, h)
		for strings.HasSuffix(path, `/index`) {
			path = strings.TrimSuffix(path, `/index`)
			a.Webx.Match(methods, path+`/`, h)
		}
	}
}

func (a *Controller) RouteByMethod() {
	if _, valid := a.Controller.(Initer); !valid {
		a.Server.Echo.Logger().Info("%T is no method Init(*Context,*App),skip.", a.Controller)
		return
	}
	t := reflect.TypeOf(a.Controller)
	e := t.Elem()
	ctlPath := e.PkgPath() + ".(*" + e.Name() + ")."
	//github.com/webx-top/{Project}/app/{App}/controller.(*Index).
	ctl := strings.ToLower(e.Name())

	for i := t.NumMethod() - 1; i >= 0; i-- {
		m := t.Method(i)
		name := m.Name
		fn := name
		h := func(u *Url) func(ctx echo.Context) error {
			return func(ctx echo.Context) error {
				c := X(ctx)
				c.Init(e.Name(), name)
				if !u.ValidExt(c.Format) {
					return c.HTML(404, `The contents can not be displayed in this format: `+c.Format)
				}
				v := reflect.New(e)
				ac := v.Interface()
				ac.(Initer).Init(c, a.App)
				if a.HasBefore {
					if err := ac.(Before).Before(); err != nil {
						return err
					}
					if c.Exit {
						return nil
					}
				}
				m := v.MethodByName(fn)
				if r, err := a.SafelyCall(m, []reflect.Value{}); err != nil {
					return err
				} else if len(r) > 0 {
					if err, ok := r[0].Interface().(error); ok && err != nil {
						return err
					}
				}
				if a.HasAfter {
					if c.Exit {
						return nil
					}
					return ac.(After).After()
				}
				return nil
			}
		}
		if strings.HasSuffix(name, `_ANY`) {
			name = strings.TrimSuffix(name, `_ANY`)
			path := "/" + ctl + "/" + strings.ToLower(name)
			u := a.App.Server.URL.SetByKey(path, ctlPath+name+"-fm")
			handler := h(u)
			a.Webx.Any(path, handler)
			for strings.HasSuffix(path, `/index`) {
				path = strings.TrimSuffix(path, `/index`)
				a.Webx.Any(path+`/`, handler)
			}
			continue
		}
		matches := methodSuffixRegex.FindAllString(name, 1)
		if len(matches) < 1 {
			continue
		}
		methods := strings.Split(strings.TrimPrefix(matches[0], `_`), `_`)
		name = strings.TrimSuffix(name, matches[0])
		path := "/" + ctl + "/" + strings.ToLower(name)
		u := a.App.Server.URL.SetByKey(path, ctlPath+name+"-fm")
		handler := h(u)
		a.Webx.Match(methods, path, handler)
		for strings.HasSuffix(path, `/index`) {
			path = strings.TrimSuffix(path, `/index`)
			a.Webx.Match(methods, path+`/`, handler)
		}
	}
}

//注册路由：Controller.AutoRoute()
func (a *Controller) AutoRoute() {
	//a.RouteByTag()
	a.RouteByMethod()
}

// safelyCall invokes `function` in recover block
func (a *Controller) SafelyCall(fn reflect.Value, args []reflect.Value) (resp []reflect.Value, err error) {
	defer func() {
		if e := recover(); e != nil {
			resp = nil
			var content string
			content = fmt.Sprintf("Handler crashed with error: %v", e)
			for i := 1; ; i += 1 {
				_, file, line, ok := runtime.Caller(i)
				if !ok {
					break
				} else {
					content += "\n"
				}
				content += fmt.Sprintf("%v %v", file, line)
			}
			a.Server.Echo.Logger().Error(content)
			err = errors.New(content)
			return
		}
	}()
	if fn.Type().NumIn() > 0 {
		return fn.Call(args), err
	}
	return fn.Call(nil), err
}
