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
package cachestore

import (
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

// Memcache adapter.
type MemcacheCache struct {
	c        *memcache.Client
	LifeTime int32
	Debug    bool
}

// create new memcache adapter.
func NewMemCache(conn []string, lifeTime int32) *MemcacheCache {
	rc := &MemcacheCache{}
	rc.c = memcache.New(conn...)
	rc.LifeTime = lifeTime
	return rc
}

// get value from memcache.
func (rc *MemcacheCache) Get(key string) (interface{}, error) {
	val, err := rc.c.Get(Md5(key))
	if err != nil || val == nil {
		if err != nil && rc.Debug {
			log.Println("[Memcache]GetErr: ", err, "Key:", key)
		}
		return nil, err
	}

	var v interface{}
	err = Decode(val.Value, &v)
	if err != nil {
		if rc.Debug {
			log.Println("[Memcache]DecodeErr: ", err, "Key:", key)
		}
		return nil, err
	}

	if rc.Debug {
		log.Println("[Memcache]Get: ", key)
	}
	return v, err
}

// put value to memcache. only support string.
func (rc *MemcacheCache) Put(key string, value interface{}) error {
	val, err := Encode(value)
	if err != nil {
		if rc.Debug {
			log.Println("[Memcache]EncodeErr: ", err, "Key:", key)
		}
		return err
	}
	item := &memcache.Item{Key: Md5(key), Value: val, Expiration: rc.LifeTime}
	err = rc.c.Set(item)
	if err != nil {
		if rc.Debug {
			log.Println("[Memcache]PutErr: ", err, "Key:", key)
		}
		return err
	}
	if rc.Debug {
		log.Println("[Memcache]Put: ", key)
	}
	return err
}

// delete value in memcache.
func (rc *MemcacheCache) Del(key string) error {
	err := rc.c.Delete(Md5(key))
	if err != nil {
		if rc.Debug {
			log.Println("[Memcache]DelErr: ", err, "Key:", key)
		}
		return err
	}
	if rc.Debug {
		log.Println("[Memcache]Del: ", key)
	}
	return err
}

// [Not Support]
// increase counter.
func (rc *MemcacheCache) Incr(key string, delta uint64) error {
	_, err := rc.c.Increment(key, delta)
	return err
}

// [Not Support]
// decrease counter.
func (rc *MemcacheCache) Decr(key string, delta uint64) error {
	_, err := rc.c.Decrement(key, delta)
	return err
}

// check value exists in memcache.
func (rc *MemcacheCache) IsExist(key string) bool {
	v, err := rc.c.Get(key)
	if err != nil || v == nil {
		return false
	}
	return true
}

// clear all cached in memcache.
func (rc *MemcacheCache) ClearAll() error {
	return rc.c.FlushAll()
}
