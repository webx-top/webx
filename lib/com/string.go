// Copyright 2013 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package com

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	r "math/rand"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// md5 hash string
func Md5(str string) string {
	m := md5.New()
	io.WriteString(m, str)
	return fmt.Sprintf("%x", m.Sum(nil))
}

func ByteMd5(b []byte) string {
	m := md5.New()
	m.Write(b)
	return hex.EncodeToString(m.Sum(nil))
}

func Token(key string, val []byte, args ...string) string {
	hm := hmac.New(sha1.New, []byte(key))
	hm.Write(val)
	for _, v := range args {
		hm.Write([]byte(v))
	}
	return fmt.Sprintf("%02x", hm.Sum(nil))
}

func Encode(data interface{}) ([]byte, error) {
	//return JsonEncode(data)
	return GobEncode(data)
}

func Decode(data []byte, to interface{}) error {
	//return JsonDecode(data, to)
	return GobDecode(data, to)
}

func GobEncode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(&data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GobDecode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

func JsonEncode(data interface{}) ([]byte, error) {
	val, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return val, nil
}

func JsonDecode(data []byte, to interface{}) error {
	return json.Unmarshal(data, to)
}

func sha(m hash.Hash, str string) string {
	io.WriteString(m, str)
	return fmt.Sprintf("%x", m.Sum(nil))
}

// sha1 hash string
func Sha1(str string) string {
	return sha(sha1.New(), str)
}

// sha256 hash string
func Sha256(str string) string {
	return sha(sha256.New(), str)
}

// trim space on left
func Ltrim(str string) string {
	return strings.TrimLeftFunc(str, unicode.IsSpace)
}

// trim space on right
func Rtrim(str string) string {
	return strings.TrimRightFunc(str, unicode.IsSpace)
}

// trim space in all string length
func Trim(str string) string {
	return strings.TrimSpace(str)
}

// repeat string times
func StrRepeat(str string, times int) string {
	return strings.Repeat(str, times)
}

// replace find all occurs to string
func StrReplace(str string, find string, to string) string {
	return strings.Replace(str, find, to, -1)
}

// IsLetter returns true if the 'l' is an English letter.
func IsLetter(l uint8) bool {
	n := (l | 0x20) - 'a'
	if n >= 0 && n < 26 {
		return true
	}
	return false
}

// Expand replaces {k} in template with match[k] or subs[atoi(k)] if k is not in match.
func Expand(template string, match map[string]string, subs ...string) string {
	var p []byte
	var i int
	for {
		i = strings.Index(template, "{")
		if i < 0 {
			break
		}
		p = append(p, template[:i]...)
		template = template[i+1:]
		i = strings.Index(template, "}")
		if s, ok := match[template[:i]]; ok {
			p = append(p, s...)
		} else {
			j, _ := strconv.Atoi(template[:i])
			if j >= len(subs) {
				p = append(p, []byte("Missing")...)
			} else {
				p = append(p, subs[j]...)
			}
		}
		template = template[i+1:]
	}
	p = append(p, template...)
	return string(p)
}

// Reverse s string, support unicode
func Reverse(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, rune := range s {
		n--
		runes[n] = rune
	}
	return string(runes[n:])
}

// RandomCreateBytes generate random []byte by specify chars.
func RandomCreateBytes(n int, alphabets ...byte) []byte {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	var randby bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randby = true
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			if randby {
				bytes[i] = alphanum[r.Intn(len(alphanum))]
			} else {
				bytes[i] = alphanum[b%byte(len(alphanum))]
			}
		} else {
			if randby {
				bytes[i] = alphabets[r.Intn(len(alphabets))]
			} else {
				bytes[i] = alphabets[b%byte(len(alphabets))]
			}
		}
	}
	return bytes
}

// Substr returns the substr from start to length.
func Substr(s string, dot string, lengthAndStart ...int) string {
	var start, length, argsLen, ln int
	argsLen = len(lengthAndStart)
	if argsLen > 0 {
		length = lengthAndStart[0]
	}
	if argsLen > 1 {
		start = lengthAndStart[1]
	}
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	ln = len(bt)
	if start > ln {
		start = start % ln
	}
	var end int = start + length
	if end > (ln - 1) {
		end = ln
	}
	if dot == "" || end == ln {
		return string(bt[start:end])
	}
	return string(bt[start:end]) + dot
}
