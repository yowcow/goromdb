package bdbstorage

import (
	"sync"
	"sync/atomic"

	"github.com/ajiyoshi-vg/goberkeleydb/bdb"
	"github.com/yowcow/goromdb/storage"
)

var (
	_ storage.Storage   = (*Storage)(nil)
	_ storage.NSStorage = (*Storage)(nil)
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

// NewNS creates and returns a storage
func NewNS() *Storage {
	return New()
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

// LoadAndIterate loads new db handle into storage, iterate through newly loaded db handle, and closes old db handle if exists
func (s *Storage) LoadAndIterate(file string, fn storage.IterationFunc) error {
	newDB, err := openBDB(file)
	if err != nil {
		return err
	}
	err = iterate(newDB, fn)
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

// GetNS finds a given ns+key in db, and returns its value
func (s *Storage) GetNS(ns, key []byte) ([]byte, error) {
	fullKey := make([]byte, 0, len(ns)+len(key))
	fullKey = append(fullKey, ns...)
	fullKey = append(fullKey, key...)
	return s.Get(fullKey)
}

func iterate(db *bdb.BerkeleyDB, fn storage.IterationFunc) error {
	cur, err := db.NewCursor(bdb.NoTxn, 0)
	if err != nil {
		return err
	}
	defer cur.Close()

	for k, v, err := cur.First(); err == nil; k, v, err = cur.Next() {
		if err = fn(k, v); err != nil {
			return err
		}
	}
	return nil
}
