package com

import (
	"fmt"
	"github.com/admpub/i18n"
)

var defaultI18n *I18n

type I18n struct {
	*i18n.TranslatorFactory
	*i18n.Translator
}

func NewI18n(rulesPath, messagesPath string) *I18n {
	f, _ := i18n.NewTranslatorFactory(
		[]string{rulesPath},
		[]string{messagesPath},
		"en",
	)
	return &I18n{
		TranslatorFactory: f,
	}
}

func (a *I18n) Get(langCode string) *i18n.Translator {
	var (
		t      *i18n.Translator
		errors []error
	)
	t, errors = a.TranslatorFactory.GetTranslator(langCode)
	_ = errors
	a.Translator = t
	return t
}

func (a *I18n) T(key string, args map[string]string) string {
	// WELCOME_MSG => "Welcome!"
	translation, _ := a.Translator.Translate(key, args)
	return translation
}

//多语言翻译
func T(key, args ...interface{}) string {
	if len(args) > 0 {
		if v, ok := arg[0].(map[string]string); ok {
			if defaultI18n == nil {
				return key
			} else {
				return defaultI18n.T(key, v)
			}
		} else {
			key = fmt.Sprintf(key, args...)
		}
	}
	if defaultI18n == nil {
		return key
	}
	return defaultI18n.T(key, map[string]string{})
}
