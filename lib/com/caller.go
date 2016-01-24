package com

import (
	"path"
	"reflect"
	"runtime"
	"strings"
)

func FuncName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

//返回包名、实例名和函数名
func FuncPath(i interface{}) (pkgName string, objName string, funcName string) {
	s := FuncName(i)
	_, file := path.Split(s)
	return ParseFuncName(file)
}

//返回完整路径包名、实例名和函数名
func FuncFullPath(i interface{}) (pkgName string, objName string, funcName string) {
	return ParseFuncName(FuncName(i))
}

func ParseFuncName(funcString string) (pkgName string, objName string, funcName string) {
	if strings.HasSuffix(funcString, `-fm`) {
		funcString = strings.TrimSuffix(funcString, `-fm`)
		part := strings.Split(funcString, `.`)
		switch len(part) {
		case 3:
			funcName = part[2]
			fallthrough
		case 2:
			objName = part[1]
			if objName[0] == '(' {
				objName = objName[1 : len(objName)-1]
				objName = strings.TrimPrefix(objName, `*`)
			}
			fallthrough
		case 1:
			pkgName = part[0]
		}
		return
	}
	part := strings.Split(funcString, `.`)
	switch len(part) {
	case 2:
		funcName = part[1]
		fallthrough
	case 1:
		pkgName = part[0]
	}
	return
}
