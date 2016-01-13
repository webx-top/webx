package main

import (
	"fmt"
	"time"

	"github.com/webx-top/webx/lib/tplex"
)

func main() {
	tpl := tplex.New("./template/")
	tpl.InitMgr(true)
	for i := 0; i < 5; i++ {
		ts := time.Now()
		fmt.Printf("==========%v: %v========\\\n", i, ts)
		tmpl := tpl.Fetch("test", nil)
		str := tpl.Parse(tmpl, map[string]interface{}{
			"test": "one---" + fmt.Sprintf("%v", i),
			"r":    []string{"one", "two", "three"},
		})
		fmt.Printf("%v\n", str)
		fmt.Printf("==========cost: %v========/\n", time.Now().Sub(ts).Seconds())
	}
}
