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
	"github.com/webx-top/webx/lib/forms/common"
	"reflect"
)

// StaticField returns a static field with the provided name and content
func StaticField(name, content string) *Field {
	ret := FieldWithType(name, formcommon.STATIC)
	ret.SetText(content)
	return ret
}

// RadioFieldFromInstance creates and initializes a radio field based on its name, the reference object instance and field number.
func StaticFieldFromInstance(val reflect.Value, t reflect.Type, fieldNo int, name string) *Field {
	ret := StaticField(name, fmt.Sprintf("%s", val.Field(fieldNo).Interface()))
	return ret
}
