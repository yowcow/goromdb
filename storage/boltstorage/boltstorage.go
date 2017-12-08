package boltstorage

import (
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
	db, err := openDB(file)
	if err != nil {
		return err
	}

	mux.Lock()
	defer mux.Unlock()

	if s.db != nil {
		s.db.Close()
	}
	s.db = db
	return nil
}

func (s *Storage) LoadAndIterate(file string, fn storage.IterationFunc, mux *sync.RWMutex) error {
	db, err := openDB(file)
	if err != nil {
		return err
	}
	err = iterate(db, s.bucket, fn)
	if err != nil {
		return err
	}

	mux.Lock()
	defer mux.Unlock()

	if s.db != nil {
		s.db.Close()
	}
	s.db = db
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

func openDB(file string) (*bolt.DB, error) {
	return bolt.Open(file, 0644, &bolt.Options{ReadOnly: true})
}

func iterate(db *bolt.DB, bucket []byte, fn storage.IterationFunc) error {
	return db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(bucket).Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			if err := fn(k, v); err != nil {
				return err
			}
		}
		return nil
	})
}
