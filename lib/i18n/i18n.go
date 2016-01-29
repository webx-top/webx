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
package i18n

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/admpub/i18n"
)

var defaultI18n *I18n

type I18n struct {
	*i18n.TranslatorFactory
	Translators map[string]*i18n.Translator
}

func New(rulesPath, messagesPath string, langCode string, defaultLangCode string) *I18n {
	f, _ := i18n.NewTranslatorFactory(
		[]string{rulesPath},
		[]string{messagesPath},
		defaultLangCode,
	)
	a := &I18n{
		TranslatorFactory: f,
		Translators:       make(map[string]*i18n.Translator),
	}
	if defaultI18n == nil {
		defaultI18n = a
	}
	a.Get(langCode)
	return a
}

func (a *I18n) Get(langCode string) *i18n.Translator {
	var (
		t      *i18n.Translator
		errors []error
	)
	t, errors = a.TranslatorFactory.GetTranslator(langCode)
	_ = errors
	a.Translators[langCode] = t
	return t
}

func (a *I18n) Reload(langCode string) {
	if strings.HasSuffix(langCode, `.yaml`) {
		langCode = strings.TrimSuffix(langCode, `.yaml`)
		langCode = filepath.Base(langCode)
	}
	a.TranslatorFactory.Reload(langCode)
	if _, ok := a.Translators[langCode]; ok {
		delete(a.Translators, langCode)
	}
}

func (a *I18n) T(langCode, key string, args map[string]string) string {
	t, ok := a.Translators[langCode]
	if !ok {
		t = a.Get(langCode)
	}
	translation, _ := t.Translate(key, args)
	return translation
}

//多语言翻译
func T(langCode, key string, args ...interface{}) string {
	if len(args) > 0 {
		if v, ok := args[0].(map[string]string); ok {
			if defaultI18n == nil {
				return key
			} else {
				return defaultI18n.T(langCode, key, v)
			}
		} else {
			key = fmt.Sprintf(key, args...)
		}
	}
	if defaultI18n == nil {
		return key
	}
	return defaultI18n.T(langCode, key, map[string]string{})
}
