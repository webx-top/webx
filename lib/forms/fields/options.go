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
package fields

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/webx-top/webx/lib/forms/common"
)

// Id - Value pair used to define an option for select and redio input fields.
type InputChoice struct {
	Id, Val string
	Checked bool
}

type ChoiceIndex struct {
	Group string
	Index int
}

func defaultValue(val reflect.Value, t reflect.Type, fieldNo int) string {
	field := val.Field(fieldNo)
	var v string = formcommon.TagVal(t, fieldNo, "form_value")
	if v == "" {
		v = fmt.Sprintf("%v", field.Interface())
	} else {
		if value := fmt.Sprintf("%v", field.Interface()); value != "" {
			v = value
		}
		/*
			switch field.Kind() {
			case reflect.Int, reflect.Int64, reflect.Int32:
				value := field.Int()
				if value != 0 {
					v = fmt.Sprintf("%v", value)
				}
			case reflect.Uint, reflect.Uint64, reflect.Uint32:
				value := field.Uint()
				if value != 0 {
					v = fmt.Sprintf("%v", value)
				}
			case reflect.Float32, reflect.Float64:
				value := field.Float()
				if value != 0 {
					v = fmt.Sprintf("%v", value)
				}
			default:
				value := fmt.Sprintf("%v", field.Interface())
				if value != "" {
					v = value
				}
			}
		*/
	}
	return v
}

// =============== RADIO

// RadioField creates a default radio button input field with the provided name and list of choices.
func RadioField(name string, choices []InputChoice) *Field {
	ret := FieldWithType(name, formcommon.RADIO)
	ret.choices = []InputChoice{}
	ret.SetChoices(choices)
	return ret
}

// RadioFieldFromInstance creates and initializes a radio field based on its name, the reference object instance and field number.
// This method looks for "form_choices" and "form_value" tags to add additional parameters to the field. "form_choices" tag is a list
// of <id>|<value> options, joined by "|" character; ex: "A|Option A|B|Option B" translates into 2 options: <A, Option A> and <B, Option B>.
func RadioFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string) *Field {
	choices := strings.Split(formcommon.TagVal(t, fieldNo, "form_choices"), "|")
	chArr := make([]InputChoice, 0)
	ret := RadioField(name, chArr)
	chMap := make(map[string]string)
	for i := 0; i < len(choices)-1; i += 2 {
		ret.choiceKeys[choices[i]] = ChoiceIndex{Group: "", Index: len(chArr)}
		chArr = append(chArr, InputChoice{choices[i], formcommon.LabelFn(choices[i+1]), false})
		chMap[choices[i]] = choices[i+1]
	}
	ret.SetChoices(chArr, false)
	v := defaultValue(val, t, fieldNo)
	if _, ok := chMap[v]; ok {
		ret.SetValue(v)
	}
	return ret
}

// ================ SELECT

// SelectField creates a default select input field with the provided name and map of choices. Choices for SelectField are grouped
// by name (if <optgroup> is needed); "" group is the default one and does not trigger a <optgroup></optgroup> rendering.
func SelectField(name string, choices map[string][]InputChoice) *Field {
	ret := FieldWithType(name, formcommon.SELECT)
	ret.choices = map[string][]InputChoice{}
	ret.SetChoices(choices)
	return ret
}

// SelectFieldFromInstance creates and initializes a select field based on its name, the reference object instance and field number.
// This method looks for "form_choices" and "form_value" tags to add additional parameters to the field. "form_choices" tag is a list
// of <group<|<id>|<value> options, joined by "|" character; ex: "G1|A|Option A|G1|B|Option B" translates into 2 options in the same group G1:
// <A, Option A> and <B, Option B>. "" group is the default one.
func SelectFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string, options map[string]struct{}) *Field {
	choices := strings.Split(formcommon.TagVal(t, fieldNo, "form_choices"), "|")
	chArr := make(map[string][]InputChoice)
	ret := SelectField(name, chArr)
	chMap := make(map[string]string)
	for i := 0; i < len(choices)-2; i += 3 {
		optgroupLabel := formcommon.LabelFn(choices[i])
		if _, ok := chArr[optgroupLabel]; !ok {
			chArr[optgroupLabel] = make([]InputChoice, 0)
		}
		id := choices[i+1]
		ret.choiceKeys[id] = ChoiceIndex{Group: optgroupLabel, Index: len(chArr[optgroupLabel])}
		chArr[optgroupLabel] = append(chArr[optgroupLabel], InputChoice{id, formcommon.LabelFn(choices[i+2]), false})
		chMap[id] = choices[i+2]
	}
	ret.SetChoices(chArr, false)
	if _, ok := options["multiple"]; ok {
		ret.MultipleChoice()
	}
	v := defaultValue(val, t, fieldNo)
	if _, ok := chMap[v]; ok {
		ret.SetValue(v)
	}
	return ret
}

// ================== CHECKBOX

func CheckboxField(name string, choices []InputChoice) *Field {
	ret := FieldWithType(name, formcommon.CHECKBOX)
	ret.choices = []InputChoice{}
	ret.SetChoices(choices)
	if len(ret.choices.([]InputChoice)) > 1 {
		ret.MultipleChoice()
	}
	return ret
}

func CheckboxFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string) *Field {
	choices := strings.Split(formcommon.TagVal(t, fieldNo, "form_choices"), "|")
	chArr := make([]InputChoice, 0)
	ret := CheckboxField(name, chArr)
	chMap := make(map[string]string)
	for i := 0; i < len(choices)-1; i += 2 {
		ret.choiceKeys[choices[i]] = ChoiceIndex{Group: "", Index: len(chArr)}
		chArr = append(chArr, InputChoice{choices[i], formcommon.LabelFn(choices[i+1]), false})
		chMap[choices[i]] = choices[i+1]
	}
	ret.SetChoices(choices, false)
	if len(ret.choices.([]InputChoice)) > 1 {
		ret.MultipleChoice()
	}
	v := defaultValue(val, t, fieldNo)
	if _, ok := chMap[v]; ok {
		ret.SetValue(v)
	}
	return ret
}

// Checkbox creates a default checkbox field with the provided name. It also makes it checked by default based
// on the checked parameter.
func Checkbox(name string, checked bool) *Field {
	ret := FieldWithType(name, formcommon.CHECKBOX)
	if checked {
		ret.AddTag("checked")
	}
	return ret
}

// CheckboxFromInstance creates and initializes a checkbox field based on its name, the reference object instance, field number and field options.
func CheckboxFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string, options map[string]struct{}) *Field {
	ret := FieldWithType(name, formcommon.CHECKBOX)
	if _, ok := options["checked"]; ok {
		ret.AddTag("checked")
	} else {
		v := val.Field(fieldNo).Bool()
		if v {
			ret.AddTag("checked")
		}
	}
	return ret
}
