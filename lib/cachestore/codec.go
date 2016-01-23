package cachestore

import (
	"github.com/webx-top/webx/lib/com"
	"log"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
}

func Md5(v string) string {
	return com.Md5(v)
}

func Encode(data interface{}) ([]byte, error) {
	return com.GobEncode(data)
}

func Decode(data []byte, to interface{}) error {
	return com.GobDecode(data, to)
}
