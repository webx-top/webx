package main

import (
	"github.com/coscms/webx"
)

type MainAction struct {
	*webx.Action

	home webx.Mapper `webx:"/"`
}

func (this *MainAction) Home() error {
	str := this.App.TemplateEx.Fetch("test", nil, map[string]interface{}{
		"test": "---one---",
		"r":    []string{"one", "two", "three"},
	})
	this.SetBody([]byte(str))
	return nil
}

func main() {
	webx.AddAction(&MainAction{})
	webx.RootApp().AppConfig.TemplateDir = `../template`
	webx.MainServer().Config.Profiler = true
	webx.Run("0.0.0.0:8888")
}
