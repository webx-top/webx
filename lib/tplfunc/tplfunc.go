/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package tplfunc

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"

	captcha "github.com/webx-top/webx/lib/captcha"
	"github.com/webx-top/webx/lib/com"
)

var TplFuncMap template.FuncMap = template.FuncMap{
	"Now":             Now,
	"Eq":              Eq,
	"Add":             Add,
	"Sub":             Sub,
	"IsNil":           IsNil,
	"Html":            Html,
	"Js":              Js,
	"Css":             Css,
	"HtmlAttr":        HtmlAttr,
	"ToHtmlAttrs":     ToHtmlAttrs,
	"ElapsedMemory":   com.ElapsedMemory,  //内存消耗
	"TotalRunTime":    com.TotalRunTime,   //运行时长(从启动服务时算起)
	"CaptchaForm":     CaptchaForm,        //验证码图片
	"FormatByte":      com.FormatByte,     //字节转为适合理解的格式
	"FormatPastTime":  com.FormatPastTime, //以前距离现在多长时间
	"DateFormat":      com.DateFormat,
	"DateFormatShort": com.DateFormatShort,
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
	"Default":         Default,
	"JsonEncode":      com.SetJson,
	"Set":             Set,
	"Append":          Append,
	"Nl2br":           NlToBr,
}

func Default(defaultV interface{}, v interface{}) interface{} {
	switch v.(type) {
	case nil:
		return defaultV
	case string:
		val, _ := v.(string)
		if val == `` {
			return defaultV
		}
	case uint8, int8, uint, int, uint32, int32, int64, uint64:
		if com.Int64(v) == 0 {
			return defaultV
		}
	case float32, float64:
		if com.Float64(v) == 0.0 {
			return defaultV
		}
	default:
		if com.Str(v) == `` {
			return defaultV
		}
	}
	return v
}

func Set(renderArgs map[string]interface{}, key string, value interface{}) string {
	renderArgs[key] = value
	return ``
}

func Append(renderArgs map[string]interface{}, key string, value interface{}) string {
	if renderArgs[key] == nil {
		renderArgs[key] = []interface{}{value}
	} else {
		renderArgs[key] = append(renderArgs[key].([]interface{}), value)
	}
	return ``
}

// Replaces newlines with <br />
func NlToBr(text string) template.HTML {
	return template.HTML(Nl2br(text))
}

// 验证码表单域
func CaptchaForm(args ...string) template.HTML {
	var id string
	if len(args) == 0 {
		id = "captcha"
	} else {
		id = args[0]
	}
	captchaId := captcha.New()
	return template.HTML(fmt.Sprintf(`<img id="`+id+`Image" src="%scaptcha/%s.png" alt="Captcha image" onclick="this.src=this.src.split('?')[0]+'?reload='+Math.random();" /><input type="hidden" name="captchaId" id="`+id+`Id" value="%s" />`, `/`, captchaId, captchaId))
}

// 验证码验证
func CaptchaVerify(captchaSolution string, r *http.Request) bool {
	captchaId := r.FormValue("captchaId")
	if !captcha.VerifyString(captchaId, captchaSolution) {
		return false
	}
	return true
}

//将换行符替换为<br />
func Nl2br(text string) string {
	return com.Nl2br(template.HTMLEscapeString(text))
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
