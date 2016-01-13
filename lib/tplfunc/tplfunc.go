package tplfunc

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/webx-top/webx/lib/com"
)

var TplFuncMap template.FuncMap = template.FuncMap{
	"Now":         Now,
	"Eq":          Eq,
	"Add":         Add,
	"Sub":         Sub,
	"IsNil":       IsNil,
	"Html":        Html,
	"Js":          Js,
	"Css":         Css,
	"HtmlAttr":    HtmlAttr,
	"ToHtmlAttrs": ToHtmlAttrs,

	"T":               T,               //译文
	"ElapsedMemory":   ElapsedMemory,   //内存消耗
	"TotalRunTime":    TotalRunTime,    //运行时长(从启动服务时算起)
	"CaptchaFormHtml": CaptchaFormHtml, //验证码图片
	"IsEmpty":         IsEmpty,         //是空(nil or len(x)==0)
	"IconHtml":        IconHtml,        //文字图标HTML
	"NotEmpty":        NotEmpty,        //非空
	"FormatByte":      FormatByte,      //字节转为适合理解的格式
	"FormatPastTime":  FormatPastTime,  //以前距离现在多长时间
	"Hook":            Hook,            //钩子
	"Icons":           Icons,           //文字图标
	"LinkJs":          LinkJs,
	"LinkCss":         LinkCss,
	"WrapHtml":        WrapHtml, //INNER文本
	"DateFormat":      DateFormat,
	"DateFormatShort": DateFormatShort,
	"Long2IP":         Interface2IP,
	"Replace":         strings.Replace, //strings.Replace(s, old, new, n)
	"Contains":        strings.Contains,
	"HasPrefix":       strings.HasPrefix,
	"HasSuffix":       strings.HasSuffix,
	"Split":           strings.Split,
	"Join":            strings.Join,
	"Str":             com.Str,
	"Int":             com.Int,
	"Int32":           com.Int32,
	"Int64":           com.Int64,
	"Float32":         com.Float32,
	"Float64":         com.Float64,
	"InSlice":         com.InSlice,
	"InSliceI":        com.InSliceIface,
	"Substr":          com.Substr,
	"StripTags":       com.StripTags,
	"Default": func(defaultV interface{}, v interface{}) interface{} {
		if v == nil || com.Str(v) == "" {
			return defaultV
		}
		return v
	},
	"JsonEncode": func(v interface{}) string {
		return Json(v, false)
	},

	"Set": func(renderArgs map[string]interface{}, key string, value interface{}) template.HTML {
		renderArgs[key] = value
		return template.HTML("")
	},
	"Append": func(renderArgs map[string]interface{}, key string, value interface{}) template.HTML {
		if renderArgs[key] == nil {
			renderArgs[key] = []interface{}{value}
		} else {
			renderArgs[key] = append(renderArgs[key].([]interface{}), value)
		}
		return template.HTML("")
	},
	// Replaces newlines with <br />
	"Nl2br": func(text string) template.HTML {
		return template.HTML(Nl2br(text))
	},
}

func IsNil(a interface{}) bool {
	switch a.(type) {
	case nil:
		return true
	}
	return false
}

func Add(left interface{}, right interface{}) interface{} {
	var rleft, rright int64
	var fleft, fright float64
	var isInt bool = true
	switch left.(type) {
	case int:
		rleft = int64(left.(int))
	case int8:
		rleft = int64(left.(int8))
	case int16:
		rleft = int64(left.(int16))
	case int32:
		rleft = int64(left.(int32))
	case int64:
		rleft = left.(int64)
	case float32:
		fleft = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	switch right.(type) {
	case int:
		rright = int64(right.(int))
	case int8:
		rright = int64(right.(int8))
	case int16:
		rright = int64(right.(int16))
	case int32:
		rright = int64(right.(int32))
	case int64:
		rright = right.(int64)
	case float32:
		fright = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	var intSum int64 = rleft + rright

	if isInt {
		return intSum
	} else {
		return fleft + fright + float64(intSum)
	}
}

func Sub(left interface{}, right interface{}) interface{} {
	var rleft, rright int64
	var fleft, fright float64
	var isInt bool = true
	switch left.(type) {
	case int:
		rleft = int64(left.(int))
	case int8:
		rleft = int64(left.(int8))
	case int16:
		rleft = int64(left.(int16))
	case int32:
		rleft = int64(left.(int32))
	case int64:
		rleft = left.(int64)
	case float32:
		fleft = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	switch right.(type) {
	case int:
		rright = int64(right.(int))
	case int8:
		rright = int64(right.(int8))
	case int16:
		rright = int64(right.(int16))
	case int32:
		rright = int64(right.(int32))
	case int64:
		rright = right.(int64)
	case float32:
		fright = float64(left.(float32))
		isInt = false
	case float64:
		fleft = left.(float64)
		isInt = false
	}

	if isInt {
		return rleft - rright
	} else {
		return fleft + float64(rleft) - (fright + float64(rright))
	}
}

func Now() time.Time {
	return time.Now()
}

func Eq(left interface{}, right interface{}) bool {
	leftIsNil := (left == nil)
	rightIsNil := (right == nil)
	if leftIsNil || rightIsNil {
		if leftIsNil && rightIsNil {
			return true
		}
		return false
	}
	return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right)
}

func Html(raw string) template.HTML {
	return template.HTML(raw)
}

func HtmlAttr(raw string) template.HTMLAttr {
	return template.HTMLAttr(raw)
}

func ToHtmlAttrs(raw map[string]interface{}) (r map[template.HTMLAttr]interface{}) {
	r = make(map[template.HTMLAttr]interface{})
	for k, v := range raw {
		r[HtmlAttr(k)] = v
	}
	return
}

func Js(raw string) template.JS {
	return template.JS(raw)
}

func Css(raw string) template.CSS {
	return template.CSS(raw)
}
