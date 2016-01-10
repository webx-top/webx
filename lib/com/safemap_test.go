package com

import (
	"testing"
)

func Test_safemap(t *testing.T) {
	bm := NewSafeMap()
	if !bm.Set("testdata", 1) {
		t.Error("set Error")
	}
	if !bm.Check("testdata") {
		t.Error("check err")
	}

	if v := bm.Get("testdata"); v.(int) != 1 {
		t.Error("get err")
	}

	bm.Delete("testdata")
	if bm.Check("testdata") {
		t.Error("delete err")
	}
}
