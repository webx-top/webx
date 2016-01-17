package main

import (
	"fmt"
	"reflect"

	"github.com/webx-top/webx/lib/tagfast"
	"github.com/webx-top/webx/lib/tagfast/example/a/b/c"
)

func main() {
	m := c.Coscms{}
	v := reflect.ValueOf(m)
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		widget := tagfast.Value(t, f, "form_widget")
		fmt.Println("widget:", widget)

		valid := tagfast.Value(t, f, "valid")
		fmt.Println("valid:", valid)

		xorm := tagfast.Value(t, f, "xorm")
		fmt.Println("xorm:", xorm)

	}
	fmt.Printf("%v", tagfast.Caches())
}
