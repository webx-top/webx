package com

import (
	"strings"
	"testing"
)

func TestFuncName(t *testing.T) {
	name := FuncName(TestFuncName)
	t.Log(name)
	if !strings.HasSuffix(name, ".TestFuncName") {
		t.Error("get func name error")
	}
}
