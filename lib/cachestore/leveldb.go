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

	"github.com/syndtr/goleveldb/leveldb"
	//"reflect"
)

// LevelDB implements CacheStore provide local machine
type LevelDB struct {
	store *leveldb.DB
	Debug bool
}

func NewLevelDB(dbfile string) *LevelDB {
	db := &LevelDB{}
	if h, err := leveldb.OpenFile(dbfile, nil); err != nil {
		panic(err)
	} else {
		db.store = h
	}
	return db
}

func (s *LevelDB) Put(key string, value interface{}) error {
	val, err := Encode(value)
	if err != nil {
		if s.Debug {
			log.Println("[LevelDB]EncodeErr: ", err, "Key:", key)
		}
		return err
	}
	err = s.store.Put([]byte(key), val, nil)
	if err != nil {
		if s.Debug {
			log.Println("[LevelDB]PutErr: ", err, "Key:", key)
		}
		return err
	}
	if s.Debug {
		log.Println("[LevelDB]Put: ", key)
	}
	return err
}

func (s *LevelDB) Get(key string) (interface{}, error) {
	data, err := s.store.Get([]byte(key), nil)
	if err != nil {
		if s.Debug {
			log.Println("[LevelDB]GetErr: ", err, "Key:", key)
		}
		return nil, err
	}
	var v interface{}
	err = Decode(data, &v)
	if err != nil {
		if s.Debug {
			log.Println("[LevelDB]DecodeErr: ", err, "Key:", key)
		}
		return nil, err
	}
	if s.Debug {
		log.Println("[LevelDB]Get: ", key, v)
	}
	return v, err
}

func (s *LevelDB) Del(key string) error {
	err := s.store.Delete([]byte(key), nil)
	if err != nil {
		if s.Debug {
			log.Println("[LevelDB]DelErr: ", err, "Key:", key)
		}
		return err
	}
	if s.Debug {
		log.Println("[LevelDB]Del: ", key)
	}
	return err
}

func (s *LevelDB) Close() error {
	return s.store.Close()
}
