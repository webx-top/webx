package session

import (
	"github.com/boltdb/bolt"
	"github.com/gorilla/sessions"
	I "github.com/webx-top/webx/lib/session/ssi"
	"github.com/yosssi/boltstore/reaper"
	"github.com/yosssi/boltstore/store"
)

var boltDB *bolt.DB

type BoltStore interface {
	Store
}

func CloseBolt() {
	if boltDB == nil {
		return
	}
	boltDB.Close()
	// Invoke a reaper which checks and removes expired sessions periodically.
	reaper.Quit(reaper.Run(boltDB, reaper.Options{}))
}

//./sessions.db
func NewBoltStore(dbFile string, options I.Options, bucketName []byte, keyPairs ...[]byte) (BoltStore, error) {
	var err error
	if boltDB == nil {
		boltDB, err = bolt.Open(dbFile, 0666, nil)
		if err != nil {
			panic(err)
		}
	}
	stor, err := store.New(boltDB, store.Config{
		SessionOptions: sessions.Options{
			Path:     options.Path,
			Domain:   options.Domain,
			MaxAge:   options.MaxAge,
			Secure:   options.Secure,
			HttpOnly: options.HttpOnly,
		},
		DBOptions: store.Options{bucketName},
	}, keyPairs...)
	if err != nil {
		return nil, err
	}
	return &boltStore{stor}, nil
}

type boltStore struct {
	*store.Store
}

func (c *boltStore) Options(options I.Options) {
	/*
		c.Store.SessionOptions = sessions.Options{
			Path:     options.Path,
			Domain:   options.Domain,
			MaxAge:   options.MaxAge,
			Secure:   options.Secure,
			HttpOnly: options.HttpOnly,
		}
	*/
}
