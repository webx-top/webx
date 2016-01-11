package com

import (
	"fmt"
	"log"
	"strconv"
)

func Int64(i interface{}) int64 {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.ParseInt(in, 10, 64)
	if err != nil {
		log.Printf("string[%s] covert int64 fail. %s", in, err)
		return 0
	}
	return out
}

func Int(i interface{}) int {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.Atoi(in)
	if err != nil {
		log.Printf("string[%s] covert int fail. %s", in, err)
		return 0
	}
	return out
}

func Int32(i interface{}) int32 {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.ParseInt(in, 10, 32)
	if err != nil {
		log.Printf("string[%s] covert int32 fail. %s", in, err)
		return 0
	}
	return int32(out)
}

func Float32(i interface{}) float32 {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.ParseFloat(in, 32)
	if err != nil {
		log.Printf("string[%s] covert float32 fail. %s", in, err)
		return 0
	}
	return float32(out)
}

func Float64(i interface{}) float64 {
	in := Str(i)
	if in == "" {
		return 0
	}
	out, err := strconv.ParseFloat(in, 64)
	if err != nil {
		log.Printf("string[%s] covert float64 fail. %s", in, err)
		return 0
	}
	return out
}

func Str(v interface{}) string {
	return fmt.Sprintf("%v", v)
}
