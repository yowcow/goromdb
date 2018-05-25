package bdbstorage

import (
	"sync"
	"sync/atomic"

	"github.com/ajiyoshi-vg/goberkeleydb/bdb"
	"github.com/yowcow/goromdb/storage"
)

var (
	_ storage.Storage = (*Storage)(nil)
)

// Storage represents a BDB storage
type Storage struct {
	db  *atomic.Value
	mux *sync.Mutex
}

// New creates and returns a storage
func New() *Storage {
	return &Storage{new(atomic.Value), new(sync.Mutex)}
}

// Load loads a new db handle into storage, and closes old db handle if exists
func (s *Storage) Load(file string) error {
	newDB, err := openBDB(file)
	if err != nil {
		return err
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	oldDB := s.getDB()
	s.db.Store(newDB)
	if oldDB != nil {
		oldDB.Close(0)
	}
	return nil
}

func (s *Storage) getDB() *bdb.BerkeleyDB {
	if ptr := s.db.Load(); ptr != nil {
		return ptr.(*bdb.BerkeleyDB)
	}
	return nil
}

func openBDB(file string) (*bdb.BerkeleyDB, error) {
	return bdb.OpenBDB(bdb.NoEnv, bdb.NoTxn, file, nil, bdb.BTree, bdb.DbReadOnly, 0)
}

// Get finds a given key in db, and returns its value
func (s *Storage) Get(key []byte) ([]byte, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	db := s.getDB()
	if db == nil {
		return nil, storage.InternalError("couldn't load db")
	}
	v, err := db.Get(bdb.NoTxn, key, 0)
	if err != nil {
		return nil, storage.KeyNotFoundError(key)
	}
	return v, nil
}
