package boltstorage

import (
	"sync"
	"sync/atomic"

	"github.com/boltdb/bolt"
	"github.com/yowcow/goromdb/storage"
)

var (
	_ storage.Storage   = (*Storage)(nil)
	_ storage.NSStorage = (*Storage)(nil)
)

// Storage represents a BoltDB storage
type Storage struct {
	db     *atomic.Value
	bucket []byte
	mux    *sync.RWMutex
}

// New creates and returns a storage
func New(b string) *Storage {
	return &Storage{new(atomic.Value), []byte(b), new(sync.RWMutex)}
}

// NewNS creates and returns a storage
func NewNS() *Storage {
	return &Storage{new(atomic.Value), nil, new(sync.RWMutex)}
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

func (s *Storage) getDB() *bolt.DB {
	if ptr := s.db.Load(); ptr != nil {
		return ptr.(*bolt.DB)
	}
	return nil
}

func openDB(file string) (*bolt.DB, error) {
	return bolt.Open(file, 0644, &bolt.Options{ReadOnly: true})
}

// Get finds a given key in db, and returns its value
func (s *Storage) Get(key []byte) ([]byte, error) {
	return s.GetNS(s.bucket, key)
}

// GetNS finds a given bucket and key in db, and returns its value
func (s *Storage) GetNS(ns, key []byte) ([]byte, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if ns == nil {
		return nil, storage.InternalError("please specify bucket")
	}
	db := s.getDB()
	if db == nil {
		return nil, storage.InternalError("couldn't load db")
	}
	var val []byte
	db.View(func(tx *bolt.Tx) error {
		val = tx.Bucket(ns).Get(key)
		return nil
	})
	if val == nil {
		return nil, storage.KeyNotFoundError(key)
	}
	return val, nil

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
