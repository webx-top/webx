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

// This package provides all the input fields logic and customization methods.
package fields

import (
	"github.com/webx-top/webx/lib/forms/common"
	"github.com/webx-top/webx/lib/forms/widgets"
	"html/template"
	"strings"
)

// Field is a generic type containing all data associated to an input field.
type Field struct {
	fieldType      string
	tmpl           string
	Widget         widgets.WidgetInterface // Public Widget field for widget customization
	name           string
	class          []string
	id             string
	params         map[string]interface{}
	css            map[string]string
	label          string
	labelClass     []string
	tag            map[string]struct{}
	value          string
	helptext       string
	errors         []string
	additionalData map[string]interface{}
	choices        interface{}
	choiceKeys     map[string]ChoiceIndex
	AppendData     map[string]interface{}
	tmplStyle      string
}

// FieldInterface defines the interface an object must implement to be used in a form. Every method returns a FieldInterface object
// to allow methods chaining.
type FieldInterface interface {
	Name() string
	Render() template.HTML
	AddClass(class string) FieldInterface
	RemoveClass(class string) FieldInterface
	AddTag(class string) FieldInterface
	RemoveTag(class string) FieldInterface
	SetId(id string) FieldInterface
	SetParam(key string, value interface{}) FieldInterface
	DeleteParam(key string) FieldInterface
	AddCss(key, value string) FieldInterface
	RemoveCss(key string) FieldInterface
	SetStyle(style string) FieldInterface
	SetLabel(label string) FieldInterface
	AddLabelClass(class string) FieldInterface
	RemoveLabelClass(class string) FieldInterface
	SetValue(value string) FieldInterface
	Disabled() FieldInterface
	Enabled() FieldInterface
	SetTmpl(tmpl string, style ...string) FieldInterface
	SetHelptext(text string) FieldInterface
	AddError(err string) FieldInterface
	MultipleChoice() FieldInterface
	SingleChoice() FieldInterface
	AddSelected(opt ...string) FieldInterface
	RemoveSelected(opt string) FieldInterface
	SetChoices(choices interface{}, saveIndex ...bool) FieldInterface
	SetText(text string) FieldInterface
	SetData(key string, value interface{})
	Data() map[string]interface{}
	String() string
}

// FieldWithType creates an empty field of the given type and identified by name.
func FieldWithType(name, t string) *Field {
	return &Field{
		fieldType:      t,
		Widget:         nil,
		name:           name,
		class:          []string{},
		id:             "",
		params:         map[string]interface{}{},
		css:            map[string]string{},
		label:          "",
		labelClass:     []string{},
		tag:            map[string]struct{}{},
		value:          "",
		helptext:       "",
		errors:         []string{},
		additionalData: map[string]interface{}{},
		choices:        nil,
		choiceKeys:     map[string]ChoiceIndex{},
		AppendData:     map[string]interface{}{},
		tmplStyle:      "",
	}
}

func (f *Field) SetTmpl(tmpl string, style ...string) FieldInterface {
	f.tmpl = tmpl
	if f.tmpl != "" && f.Widget != nil {
		var s string
		if len(style) > 0 {
			s = style[0]
		} else {
			s = f.tmplStyle
		}
		f.Widget = widgets.BaseWidget(s, f.fieldType, f.tmpl)
	}
	return f
}

// SetStyle sets the style (e.g.: BASE, BOOTSTRAP) of the field, correctly populating the Widget field.
func (f *Field) SetStyle(style string) FieldInterface {
	f.tmplStyle = style
	f.Widget = widgets.BaseWidget(style, f.fieldType, f.tmpl)
	return f
}

func (f *Field) SetData(key string, value interface{}) {
	f.AppendData[key] = value
}

// Name returns the name of the field.
func (f *Field) Name() string {
	return strings.TrimSuffix(f.name, "[]")
}

func (f *Field) Data() map[string]interface{} {
	safeParams := make(map[template.HTMLAttr]interface{})
	for k, v := range f.params {
		safeParams[template.HTMLAttr(k)] = v
	}
	data := map[string]interface{}{
		"classes":      f.class,
		"id":           f.id,
		"name":         f.name,
		"params":       safeParams,
		"css":          f.css,
		"type":         f.fieldType,
		"label":        f.label,
		"labelClasses": f.labelClass,
		"tags":         f.tag,
		"value":        f.value,
		"helptext":     f.helptext,
		"errors":       f.errors,
		"container":    "form",
		"choices":      f.choices,
	}
	for k, v := range f.additionalData {
		data[k] = v
	}
	for k, v := range f.AppendData {
		data[k] = v
	}
	return data
}

// Render packs all data and executes widget render method.
func (f *Field) Render() template.HTML {
	if f.Widget != nil {
		return template.HTML(f.Widget.Render(f.Data()))
	}
	return template.HTML("")
}

func (f *Field) String() string {
	if f.Widget != nil {
		return f.Widget.Render(f.Data())
	}
	return ""
}

// AddClass adds a class to the field.
func (f *Field) AddClass(class string) FieldInterface {
	f.class = append(f.class, class)
	return f
}

// RemoveClass removes a class from the field, if it was present.
func (f *Field) RemoveClass(class string) FieldInterface {
	ind := -1
	for i, v := range f.class {
		if v == class {
			ind = i
			break
		}
	}

	if ind != -1 {
		f.class = append(f.class[:ind], f.class[ind+1:]...)
	}
	return f
}

// SetId associates the given id to the field, overwriting any previous id.
func (f *Field) SetId(id string) FieldInterface {
	f.id = id
	return f
}

// SetLabel saves the label to be rendered along with the field.
func (f *Field) SetLabel(label string) FieldInterface {
	f.label = label
	return f
}

// SetLablClass allows to define custom classes for the label.
func (f *Field) AddLabelClass(class string) FieldInterface {
	f.labelClass = append(f.labelClass, class)
	return f
}

// RemoveLabelClass removes the given class from the field label.
func (f *Field) RemoveLabelClass(class string) FieldInterface {
	ind := -1
	for i, v := range f.labelClass {
		if v == class {
			ind = i
			break
		}
	}

	if ind != -1 {
		f.labelClass = append(f.labelClass[:ind], f.labelClass[ind+1:]...)
	}
	return f
}

// SetParam adds a parameter (defined as key-value pair) in the field.
func (f *Field) SetParam(key string, value interface{}) FieldInterface {
	f.params[key] = value
	return f
}

// DeleteParam removes a parameter identified by key from the field.
func (f *Field) DeleteParam(key string) FieldInterface {
	delete(f.params, key)
	return f
}

// AddCss adds a custom CSS style the field.
func (f *Field) AddCss(key, value string) FieldInterface {
	f.css[key] = value
	return f
}

// RemoveCss removes CSS options identified by key from the field.
func (f *Field) RemoveCss(key string) FieldInterface {
	delete(f.css, key)
	return f
}

// Disabled add the "disabled" tag to the field, making it unresponsive in some environments (e.g. Bootstrap).
func (f *Field) Disabled() FieldInterface {
	f.AddTag("disabled")
	return f
}

// Enabled removes the "disabled" tag from the field, making it responsive.
func (f *Field) Enabled() FieldInterface {
	f.RemoveTag("disabled")
	return f
}

// AddTag adds a no-value parameter (e.g.: checked, disabled) to the field.
func (f *Field) AddTag(tag string) FieldInterface {
	f.tag[tag] = struct{}{}
	return f
}

// RemoveTag removes a no-value parameter from the field.
func (f *Field) RemoveTag(tag string) FieldInterface {
	delete(f.tag, tag)
	return f
}

// SetValue sets the value parameter for the field.
func (f *Field) SetValue(value string) FieldInterface {
	f.value = value
	f.AddSelected(f.value)
	return f
}

// SetHelptext saves the field helptext.
func (f *Field) SetHelptext(text string) FieldInterface {
	f.helptext = text
	return f
}

// AddError adds an error string to the field. It's valid only for Bootstrap forms.
func (f *Field) AddError(err string) FieldInterface {
	f.errors = append(f.errors, err)
	return f
}

// MultipleChoice configures the SelectField to accept and display multiple choices.
// It has no effect if type is not SELECT.
func (f *Field) MultipleChoice() FieldInterface {
	switch f.fieldType {
	case formcommon.SELECT:
		f.AddTag("multiple")
		fallthrough
	case formcommon.CHECKBOX:
		// fix name if necessary
		if !strings.HasSuffix(f.name, "[]") {
			f.name = f.name + "[]"
		}
	}
	return f
}

// SingleChoice configures the Field to accept and display only one choice (valid for SelectFields only).
// It has no effect if type is not SELECT.
func (f *Field) SingleChoice() FieldInterface {
	switch f.fieldType {
	case formcommon.SELECT:
		f.RemoveTag("multiple")
		fallthrough
	case formcommon.CHECKBOX:
		if strings.HasSuffix(f.name, "[]") {
			f.name = strings.TrimSuffix(f.name, "[]")
		}
	}
	return f
}

// If the field is configured as "multiple", AddSelected adds a selected value to the field (valid for SelectFields only).
// It has no effect if type is not SELECT.
func (f *Field) AddSelected(opt ...string) FieldInterface {
	switch f.fieldType {
	case formcommon.SELECT:
		for _, v := range opt {
			i := f.choiceKeys[v]
			if vc, ok := f.choices.(map[string][]InputChoice)[i.Group]; ok {
				if len(vc) > i.Index {
					f.choices.(map[string][]InputChoice)[i.Group][i.Index].Checked = true
				}
			}
		}
	case formcommon.RADIO, formcommon.CHECKBOX:
		size := len(f.choices.([]InputChoice))
		for _, v := range opt {
			i := f.choiceKeys[v]
			if size > i.Index {
				f.choices.([]InputChoice)[i.Index].Checked = true
			}
		}
	}
	return f
}

// If the field is configured as "multiple", AddSelected removes the selected value from the field (valid for SelectFields only).
// It has no effect if type is not SELECT.
func (f *Field) RemoveSelected(opt string) FieldInterface {
	switch f.fieldType {
	case formcommon.SELECT:
		i := f.choiceKeys[opt]
		if vc, ok := f.choices.(map[string][]InputChoice)[i.Group]; ok {
			if len(vc) > i.Index {
				f.choices.(map[string][]InputChoice)[i.Group][i.Index].Checked = false
			}
		}

	case formcommon.RADIO, formcommon.CHECKBOX:
		size := len(f.choices.([]InputChoice))
		i := f.choiceKeys[opt]
		if size > i.Index {
			f.choices.([]InputChoice)[i.Index].Checked = false
		}
	}
	return f
}

// SetChoices takes as input a dictionary whose key-value entries are defined as follows: key is the group name (the empty string
// is the default group that is not explicitly rendered) and value is the list of choices belonging to that group.
// Grouping is only useful for Select fields, while groups are ignored in Radio fields.
// It has no effect if type is not SELECT.
func (f *Field) SetChoices(choices interface{}, saveIndex ...bool) FieldInterface {
	if choices == nil {
		return f
	}
	switch f.fieldType {
	case formcommon.SELECT:
		var ch map[string][]InputChoice
		if c, ok := choices.(map[string][]InputChoice); ok {
			ch = c
		} else {
			c, _ := choices.([]InputChoice)
			ch = map[string][]InputChoice{"": c}
		}
		f.choices = ch
		if len(saveIndex) < 1 || saveIndex[0] {
			for k, v := range ch {
				for idx, ipt := range v {
					f.choiceKeys[ipt.Id] = ChoiceIndex{Group: k, Index: idx}
				}
			}
		}

	case formcommon.RADIO, formcommon.CHECKBOX:
		ch, _ := choices.([]InputChoice)
		f.choices = ch
		if len(saveIndex) < 1 || saveIndex[0] {
			for idx, ipt := range ch {
				f.choiceKeys[ipt.Id] = ChoiceIndex{Group: "", Index: idx}
			}
		}
	}
	return f
}

// SetText saves the provided text as content of the field, usually a TextAreaField.
func (f *Field) SetText(text string) FieldInterface {
	if f.fieldType == formcommon.BUTTON ||
		f.fieldType == formcommon.SUBMIT ||
		f.fieldType == formcommon.RESET ||
		f.fieldType == formcommon.STATIC ||
		f.fieldType == formcommon.TEXTAREA {
		f.additionalData["text"] = text
	}
	return f
}
