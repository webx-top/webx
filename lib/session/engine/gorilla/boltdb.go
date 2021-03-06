package session

import (
	"github.com/admpub/boltstore/reaper"
	"github.com/admpub/boltstore/store"
	"github.com/admpub/sessions"
	"github.com/boltdb/bolt"
	"github.com/webx-top/webx/lib/events"
	I "github.com/webx-top/webx/lib/session/ssi"
)

var boltDB *bolt.DB
var onCloseBolt func() error

type BoltStore interface {
	Store
}

func CloseBolt() {
	if boltDB == nil {
		return
	}
	boltDB.Close()
	if onCloseBolt != nil {
		onCloseBolt()
	}
}

//./sessions.db
func NewBoltStore(dbFile string, options I.Options, bucketName []byte, keyPairs ...[]byte) (BoltStore, error) {
	var err error
	if boltDB == nil {
		boltDB, err = bolt.Open(dbFile, 0666, nil)
		if err != nil {
			panic(err)
		}
		quiteC, doneC := reaper.Run(boltDB, reaper.Options{})
		onCloseBolt = func() error {
			// Invoke a reaper which checks and removes expired sessions periodically.
			reaper.Quit(quiteC, doneC)
			return nil
		}
		events.AddEvent(`webx.serverExit`, func(_ interface{}, next func(bool)) {
			CloseBolt()
			next(true)
		})
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
