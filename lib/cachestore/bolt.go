package cachestore

import (
	"log"

	bolt "github.com/boltdb/bolt"
)

// Bolt implements ds.Datastore
// TODO: use buckets to represent the heirarchy of the ds.Keys
type Bolt struct {
	db         *bolt.DB
	bucketName []byte
	Debug      bool
}

func NewBolt(dbFile, bucket string) (*Bolt, error) {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		return nil, err
	}

	// TODO: need to do db.Close() sometime...
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucket))
		return err
	})
	if err != nil {
		return nil, err
	}

	return &Bolt{
		db:         db,
		bucketName: []byte(bucket),
	}, nil
}

func (bd *Bolt) Close() error {
	return bd.db.Close()
}

func (bd *Bolt) Del(key string) error {
	return bd.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bd.bucketName).Delete([]byte(key))
	})
}

func (bd *Bolt) Get(key string) (interface{}, error) {
	var out []byte
	err := bd.db.View(func(tx *bolt.Tx) error {
		mmval := tx.Bucket(bd.bucketName).Get([]byte(key))
		if mmval == nil {
			return nil
		}
		out = make([]byte, len(mmval))
		copy(out, mmval)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if out == nil {
		return nil, err
	}
	var v interface{}
	err = Decode(out, &v)
	if err != nil {
		if bd.Debug {
			log.Println("[Bolt]DecodeErr: ", err, "Key:", key)
		}
		return nil, err
	}
	return v, err
}

func (bd *Bolt) ConsumeValue(key string, f func([]byte) error) error {
	return bd.db.View(func(tx *bolt.Tx) error {
		mmval := tx.Bucket(bd.bucketName).Get([]byte(key))
		if mmval == nil {
			return nil
		}
		return f(mmval)
	})
}

func (bd *Bolt) Has(key string) (bool, error) {
	var found bool
	err := bd.db.View(func(tx *bolt.Tx) error {
		val := tx.Bucket(bd.bucketName).Get([]byte(key))
		found = (val != nil)
		return nil
	})
	return found, err
}

func (bd *Bolt) Put(key string, val interface{}) error {
	bval, err := Encode(val)
	if err != nil {
		if bd.Debug {
			log.Println("[Bolt]EncodeErr: ", err, "Key:", key)
		}
		return err
	}
	return bd.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bd.bucketName).Put([]byte(key), bval)
	})
}

func (bd *Bolt) PutMany(data map[string]interface{}) error {
	return bd.db.Update(func(tx *bolt.Tx) error {
		buck := tx.Bucket(bd.bucketName)
		for k, v := range data {
			bval, err := Encode(v)
			if err != nil {
				if bd.Debug {
					log.Println("[Bolt]EncodeErr: ", err, "Key:", k)
				}
				return err
			}
			err := buck.Put([]byte(k), bval)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
