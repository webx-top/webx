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
