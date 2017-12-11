package boltstorage

import (
	"sync"

	"github.com/boltdb/bolt"
	"github.com/yowcow/goromdb/storage"
)

type Storage struct {
	db     *bolt.DB
	bucket []byte
	mux    *sync.RWMutex
}

func New(b string) *Storage {
	return &Storage{
		nil,
		[]byte(b),
		new(sync.RWMutex),
	}
}

func (s *Storage) Load(file string) error {
	newDB, err := openDB(file)
	if err != nil {
		return err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	if oldDB := s.getDB(); oldDB != nil {
		oldDB.Close()
	}
	s.db = newDB
	return nil
}

func (s *Storage) LoadAndIterate(file string, fn storage.IterationFunc) error {
	newDB, err := openDB(file)
	if err != nil {
		return err
	}
	err = iterate(newDB, s.bucket, fn)
	if err != nil {
		return err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	if oldDB := s.getDB(); oldDB != nil {
		oldDB.Close()
	}
	s.db = newDB
	return nil
}

func (s Storage) getDB() *bolt.DB {
	return s.db
}

func openDB(file string) (*bolt.DB, error) {
	return bolt.Open(file, 0644, &bolt.Options{ReadOnly: true})
}

func (s Storage) Get(key []byte) ([]byte, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	if db := s.getDB(); db != nil {
		var val []byte
		db.View(func(tx *bolt.Tx) error {
			val = tx.Bucket(s.bucket).Get(key)
			return nil
		})
		if val == nil {
			return nil, storage.KeyNotFoundError(key)
		}
		return val, nil
	}
	return nil, storage.KeyNotFoundError(key)
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
