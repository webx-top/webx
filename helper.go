package webx

import (
	"github.com/webx-top/echo"
)

func X(c echo.Context) *Context {
	return c.(*Context)
}

func MustString(c echo.Context, k string) (r string) {
	if v, ok := c.Get(k).(string); ok {
		r = v
	}
	return
}

func MustBool(c echo.Context, k string) (r bool) {
	if v, ok := c.Get(k).(bool); ok {
		r = v
	}
	return
}

func MustInt(c echo.Context, k string) (r int) {
	if v, ok := c.Get(k).(int); ok {
		r = v
	}
	return
}

func MustInt32(c echo.Context, k string) (r int32) {
	if v, ok := c.Get(k).(int32); ok {
		r = v
	}
	return
}

func MustInt64(c echo.Context, k string) (r int64) {
	if v, ok := c.Get(k).(int64); ok {
		r = v
	}
	return
}

func MustUint(c echo.Context, k string) (r uint) {
	if v, ok := c.Get(k).(uint); ok {
		r = v
	}
	return
}

func MustUint32(c echo.Context, k string) (r uint32) {
	if v, ok := c.Get(k).(uint32); ok {
		r = v
	}
	return
}

func MustUint64(c echo.Context, k string) (r uint64) {
	if v, ok := c.Get(k).(uint64); ok {
		r = v
	}
	return
}

func MustFloat32(c echo.Context, k string) (r float32) {
	if v, ok := c.Get(k).(float32); ok {
		r = v
	}
	return
}

func MustFloat64(c echo.Context, k string) (r float64) {
	if v, ok := c.Get(k).(float64); ok {
		r = v
	}
	return
}

func MustUint8(c echo.Context, k string) (r uint8) {
	if v, ok := c.Get(k).(uint8); ok {
		r = v
	}
	return
}
