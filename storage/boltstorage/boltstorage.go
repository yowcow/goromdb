package boltstorage

import (
	"sync"
	"sync/atomic"

	"github.com/boltdb/bolt"
	"github.com/yowcow/goromdb/storage"
)

// Storage represents a BoltDB storage
type Storage struct {
	db     *atomic.Value
	bucket []byte
	mux    *sync.RWMutex
}

// New creates and returns a storage
func New(b string) *Storage {
	return &Storage{
		new(atomic.Value),
		[]byte(b),
		new(sync.RWMutex),
	}
}

// Load loads a new db handle into storage, and closes old db handle if exists
func (s *Storage) Load(file string) error {
	newDB, err := openDB(file)
	if err != nil {
		return err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldDB := s.getDB()
	s.db.Store(newDB)
	if oldDB != nil {
		oldDB.Close()
	}
	return nil
}

// LoadAndIterate loads a new db handle, iterate through newly loaded db, and closes old db handle if exists
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

	oldDB := s.getDB()
	s.db.Store(newDB)
	if oldDB != nil {
		oldDB.Close()
	}
	return nil
}

func (s Storage) getDB() *bolt.DB {
	if ptr := s.db.Load(); ptr != nil {
		return ptr.(*bolt.DB)
	}
	return nil
}

func openDB(file string) (*bolt.DB, error) {
	return bolt.Open(file, 0644, &bolt.Options{ReadOnly: true})
}

// Get finds a given key in db, and returns its value
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
