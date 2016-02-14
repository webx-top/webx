package tplex

import (
	"io"
)

type TemplateEx interface {
	Init(...bool)
	SetFuncMapFn(func() map[string]interface{})
	Render(io.Writer, string, interface{}, map[string]interface{}) error
	Fetch(string, interface{}, map[string]interface{}) string
	RawContent(string) ([]byte, error)
	MonitorEvent(func(string))
	ClearCache()
	Close()
}

var engines = make(map[string]func(string) TemplateEx)

func Create(key string, tmplDir string) TemplateEx {
	if fn, ok := engines[key]; ok {
		return fn(tmplDir)
	}
	return New(tmplDir)
}

func Reg(key string, val func(string) TemplateEx) {
	engines[key] = val
}

func Del(key string) {
	if _, ok := engines[key]; ok {
		delete(engines, key)
	}
}
