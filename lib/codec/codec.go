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
package codec

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"github.com/webx-top/webx/lib/com"
	"log"
)

type Codec interface {
	Encode(rawData, authKey string) string
	Decode(cryptedData, authKey string) string
}

type AesCrypto struct {
	key map[string][]byte
}

var Default Codec = &AesCrypto{key: make(map[string][]byte, 0)}

const (
	aesKeyLen = 128
	keyLen    = aesKeyLen / 8
)

func (c *AesCrypto) aesKey(key []byte) []byte {
	if c.key == nil {
		c.key = make(map[string][]byte, 0)
	}
	ckey := string(key)
	k, ok := c.key[ckey]
	if !ok {
		if len(key) == keyLen {
			return key
		}

		k = make([]byte, keyLen)
		copy(k, key)
		for i := keyLen; i < len(key); {
			for j := 0; j < keyLen && i < len(key); j, i = j+1, i+1 {
				k[j] ^= key[i]
			}
		}
		c.key[ckey] = k
	}
	return k
}

func (c *AesCrypto) Encode(rawData, authKey string) string {
	in := []byte(rawData)
	key := []byte(authKey)
	key = c.aesKey(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
		return ""
	}
	blockSize := block.BlockSize()
	in = PKCS5Padding(in, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(in))
	blockMode.CryptBlocks(crypted, in)
	return com.Base64Encode(string(crypted))
}

func (c *AesCrypto) Decode(cryptedData, authKey string) string {
	cryptedData, _ = com.Base64Decode(cryptedData)
	if cryptedData == "" {
		return ""
	}
	in := []byte(cryptedData)
	key := []byte(authKey)
	key = c.aesKey(key)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
		return ""
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(in))
	blockMode.CryptBlocks(origData, in)
	origData = PKCS5UnPadding(origData)
	return string(origData)
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	if length > unpadding {
		return origData[:(length - unpadding)]
	}
	return []byte{}
}
