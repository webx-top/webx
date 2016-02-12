package tplex

import (
	"html/template"
	"io"
)

type TemplateEx interface {
	Init(...bool)
	SetFuncMapFn(func() template.FuncMap)
	Render(io.Writer, string, interface{}, template.FuncMap) error
	Fetch(string, interface{}, template.FuncMap) string
	RawContent(string) ([]byte, error)
	MonitorEvent(func(string))
	ClearCache()
	Close()
}
