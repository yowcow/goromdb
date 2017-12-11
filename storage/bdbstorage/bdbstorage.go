package bdbstorage

import (
	"sync"
	"sync/atomic"

	"github.com/ajiyoshi-vg/goberkeleydb/bdb"
	"github.com/yowcow/goromdb/storage"
)

type Storage struct {
	mux *sync.Mutex
	db  *atomic.Value
}

func New() *Storage {
	return &Storage{
		new(sync.Mutex),
		new(atomic.Value),
	}
}

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

func (s Storage) getDB() *bdb.BerkeleyDB {
	if ptr := s.db.Load(); ptr != nil {
		return ptr.(*bdb.BerkeleyDB)
	}
	return nil
}

func openBDB(file string) (*bdb.BerkeleyDB, error) {
	return bdb.OpenBDB(bdb.NoEnv, bdb.NoTxn, file, nil, bdb.BTree, bdb.DbReadOnly, 0)
}

func (s Storage) Get(key []byte) ([]byte, error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if db := s.getDB(); db != nil {
		v, err := db.Get(bdb.NoTxn, key, 0)
		if err != nil {
			return nil, storage.KeyNotFoundError(key)
		}
		return v, nil
	}
	return nil, storage.KeyNotFoundError(key)
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
