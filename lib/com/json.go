package com

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Read json data, writes in struct f
func GetJson(dat *string, s interface{}) {
	err := json.Unmarshal([]byte(*dat), s)
	if err != nil {
		panic("Get json failed")
	}
}

// Struct s will be converted to json format
func SetJson(s interface{}) string {
	dat, err := json.Marshal(s)
	if err != nil {
		panic("Set json failed")
	}
	return string(dat)
}

// Json data read from the specified file
func ReadJson(path string, s interface{}) {
	dat, err1 := ioutil.ReadFile(path)
	if err1 != nil {
		panic("Json file fails to open")
	}
	err2 := json.Unmarshal(dat, s)
	if err2 != nil {
		panic("Create json failed")
	}
}

// The json data is written to the specified file
func WriteJson(path string, dat *string) {
	_, err0 := os.Stat(path)
	if err0 != nil || !os.IsExist(err0) {
		os.Create(path)
	}
	err := ioutil.WriteFile(path, []byte(*dat), 0644)
	if err != nil {
		panic("Create json file failed")
	}
}

//输出对象和数组的结构信息
func Dump(m interface{}) string {
	v, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	return string(v)
}
