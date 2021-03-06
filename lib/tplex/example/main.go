package main

import (
	"fmt"
	"time"

	"github.com/webx-top/webx/lib/tplex"
)

func main() {
	tpl := tplex.New("./template/")
	tpl.Init(true)
	for i := 0; i < 5; i++ {
		ts := time.Now()
		fmt.Printf("==========%v: %v========\\\n", i, ts)
		str := tpl.Fetch("test", map[string]interface{}{
			"test": "one---" + fmt.Sprintf("%v", i),
			"r":    []string{"one", "two", "three"},
		}, nil)
		fmt.Printf("%v\n", str)
		fmt.Printf("==========cost: %v========/\n", time.Now().Sub(ts).Seconds())
	}
}
