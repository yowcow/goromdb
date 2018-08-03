package boltstorage

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/boltdb/bolt"
	"github.com/yowcow/goromdb/storage"
)

var (
	_ storage.Storage = (*Storage)(nil)
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
	s.mux.RLock()
	defer s.mux.RUnlock()

	db := s.getDB()
	if db == nil {
		return nil, storage.InternalError("couldn't load db")
	}

	return getFromBucket(db, s.bucket, key)
}

func getFromBucket(db *bolt.DB, bucket, key []byte) ([]byte, error) {
	var retVal []byte

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return storage.InternalError(fmt.Sprintf("bucket %v not found", bucket))
		}

		val := b.Get(key)
		if val == nil {
			return storage.KeyNotFoundError(key)
		}

		// Making sure that []byte is safe.
		// Without copy, returning []byte may be corrupted at the time of reference later on.
		retVal = make([]byte, len(val))
		copy(retVal, val)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return retVal, nil
}
