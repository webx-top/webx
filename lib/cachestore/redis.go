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
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/webx-top/webx/lib/cachestore/redigo/redis"
)

var (
	// the collection name of redis for cache adapter.
	DefaultKey string = "WebxRedis"
)

// Redis cache adapter.
type Redis struct {
	p        *redis.Pool // redis connection pool
	conninfo string
	dbnum    int
	key      string
	Debug    bool
	LifeTime int32
}

// create new redis cache with default collection name.
func NewRedis(cf map[string]string, lifeTime int32) *Redis {
	rc := &Redis{key: DefaultKey}
	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}
	rc.key = cf["key"]
	rc.conninfo = cf["conn"]
	rc.LifeTime = lifeTime
	rc.connectInit()
	return rc
}

// actually do the redis cmds
func (rc *Redis) do(commandName string, args ...interface{}) (reply interface{}, err error) {
	if rc.p == nil {
		rc.connectInit()
	}
	c := rc.p.Get()
	defer c.Close()

	return c.Do(commandName, args...)
}

// Get cache from redis.
func (rc *Redis) Get(key string) (interface{}, error) {
	val, err := rc.do("GET", key)
	if err != nil {
		if rc.Debug {
			log.Println("[Redis]GetErr: ", err, "Key:", key)
		}
		return nil, err
	}
	var v interface{}
	err = Decode(val.([]byte), &v)
	if err != nil {
		if rc.Debug {
			log.Println("[Redis]DecodeErr: ", err, "Key:", key)
		}
		return nil, err
	}
	if rc.Debug {
		log.Println("[Redis]Get: ", key)
	}
	return v, err
}

// put cache to redis.
// timeout is ignored.
func (rc *Redis) Put(key string, value interface{}) error {
	val, err := Encode(value)
	if err != nil {
		if rc.Debug {
			log.Println("[Redis]EncodeErr: ", err, "Key:", key)
		}
		return err
	}
	if _, err = rc.do("SETEX", key, rc.LifeTime, val); err != nil {
		return err
	}
	_, err = rc.do("HSET", rc.key, key, true)
	if err != nil {
		if rc.Debug {
			log.Println("[Redis]PutErr: ", err, "Key:", key)
		}
		return err
	}
	if rc.Debug {
		log.Println("[Redis]Put: ", key)
	}
	return err
}

// delete cache in redis.
func (rc *Redis) Del(key string) error {
	var err error
	if _, err = rc.do("DEL", key); err != nil {
		return err
	}
	_, err = rc.do("HDEL", rc.key, key)
	if err != nil {
		if rc.Debug {
			log.Println("[Redis]DelErr: ", err, "Key:", key)
		}
		return err
	}
	if rc.Debug {
		log.Println("[Redis]Del: ", key)
	}
	return err
}

// check cache exist in redis.
func (rc *Redis) IsExist(key string) bool {
	v, err := redis.Bool(rc.do("EXISTS", key))
	if err != nil {
		return false
	}
	if v == false {
		if _, err = rc.do("HDEL", rc.key, key); err != nil {
			return false
		}
	}
	return v
}

// increase counter in redis.
func (rc *Redis) Incr(key string, delta uint64) error {
	_, err := redis.Bool(rc.do("INCRBY", key, delta))
	return err
}

// decrease counter in redis.
func (rc *Redis) Decr(key string, delta uint64) error {
	_, err := redis.Bool(rc.do("INCRBY", key, delta))
	return err
}

// clean all cache in redis. delete this redis collection.
func (rc *Redis) ClearAll() error {
	cachedKeys, err := redis.Strings(rc.do("HKEYS", rc.key))
	if err != nil {
		return err
	}
	for _, str := range cachedKeys {
		if _, err = rc.do("DEL", str); err != nil {
			return err
		}
	}
	_, err = rc.do("DEL", rc.key)
	return err
}

// start redis cache adapter.
// config is like {"key":"collection key","conn":"connection info","dbnum":"0"}
// the cache item in redis are stored forever,
// so no gc operation.
func (rc *Redis) Connect(config string) error {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	if _, ok := cf["key"]; !ok {
		cf["key"] = DefaultKey
	}
	if _, ok := cf["conn"]; !ok {
		return errors.New("config has no conn key")
	}
	if _, ok := cf["dbnum"]; !ok {
		cf["dbnum"] = "0"
	}
	rc.key = cf["key"]
	rc.conninfo = cf["conn"]
	rc.dbnum, _ = strconv.Atoi(cf["dbnum"])
	rc.connectInit()

	c := rc.p.Get()
	defer c.Close()

	return c.Err()
}

// connect to redis.
func (rc *Redis) connectInit() {
	dialFunc := func() (c redis.Conn, err error) {
		c, err = redis.Dial("tcp", rc.conninfo)
		_, selecterr := c.Do("SELECT", rc.dbnum)
		if selecterr != nil {
			c.Close()
			return nil, selecterr
		}
		return
	}
	// initialize a new pool
	rc.p = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,
		Dial:        dialFunc,
	}
}
