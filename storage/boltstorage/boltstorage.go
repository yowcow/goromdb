package boltstorage

import (
	"fmt"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/yowcow/goromdb/storage"
)

type Storage struct {
	db     *bolt.DB
	bucket []byte
}

func New(b string) *Storage {
	return &Storage{nil, []byte(b)}
}

func (s *Storage) Load(file string, mux *sync.RWMutex) error {
	db, err := bolt.Open(file, 0644, &bolt.Options{ReadOnly: true})
	if err != nil {
		return err
	}

	mux.Lock()
	defer mux.Unlock()
	oldDB := s.db
	s.db = db
	if oldDB != nil {
		if err = oldDB.Close(); err != nil {
			return err
		}
		oldDB = nil
	}
	return nil
}

func (s Storage) Get(key []byte) ([]byte, error) {
	var val []byte
	s.db.View(func(tx *bolt.Tx) error {
		val = tx.Bucket(s.bucket).Get(key)
		return nil
	})
	if val == nil {
		return nil, storage.KeyNotFoundError(key)
	}
	return val, nil
}

func (s Storage) Iterate(fn storage.IteratorFunc) error {
	if s.db == nil {
		return fmt.Errorf("no boltdb handle in storage")
	}
	return s.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(s.bucket).Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if err := fn(k, v); err != nil {
				return err
			}
		}
		return nil
	})
}
