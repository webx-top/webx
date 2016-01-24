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

var methodSuffixRegex = regexp.MustCompile(`(?:_(?:` + strings.Join(echo.Methods(), `|`) + `))+$`)

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

//注册路由：Controller.AutoRoute()
func (a *Controller) AutoRoute() {
	if _, valid := a.Controller.(Initer); !valid {
		a.Server.Echo.Logger().Info("%T is no method Init(*Context,*App),skip.", a.Controller)
		return
	}
	t := reflect.TypeOf(a.Controller)
	e := t.Elem()
	ctlPath := e.PkgPath() + ".(*" + e.Name() + ")." //github.com/webx-top/{Project}/app/{App}/controller.(*Index).
	ctl := strings.ToLower(e.Name())
	for i := t.NumMethod() - 1; i >= 0; i-- {
		m := t.Method(i)
		name := m.Name
		fn := name
		h := func(ctx echo.Context) error {
			c := X(ctx)
			c.Init(e.Name(), name)
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

			if r, err := a.SafelyCall(v, fn, []reflect.Value{}); err != nil {
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
		if strings.HasSuffix(name, `_ANY`) {
			name = strings.TrimSuffix(name, `_ANY`)
			path := "/" + ctl + "/" + strings.ToLower(name)
			a.App.Server.URL.SetByKey(path, ctlPath+name+"-fm")
			a.Webx.Any(path, h)
			for strings.HasSuffix(path, `/index`) {
				path = strings.TrimSuffix(path, `/index`)
				a.Webx.Any(path+`/`, h)
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
		a.App.Server.URL.SetByKey(path, ctlPath+name+"-fm")
		a.Webx.Match(methods, path, h)
		for strings.HasSuffix(path, `/index`) {
			path = strings.TrimSuffix(path, `/index`)
			a.Webx.Match(methods, path+`/`, h)
		}
	}
	t = t.Elem()
}

// safelyCall invokes `function` in recover block
func (a *Controller) SafelyCall(vc reflect.Value, method string, args []reflect.Value) (resp []reflect.Value, err error) {
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
	fn := vc.MethodByName(method)
	if fn.Type().NumIn() > 0 {
		return fn.Call(args), err
	}
	return fn.Call(nil), err
}
