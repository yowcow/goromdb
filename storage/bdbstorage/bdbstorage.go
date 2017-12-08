package bdbstorage

import (
	"sync"

	"github.com/ajiyoshi-vg/goberkeleydb/bdb"
	"github.com/yowcow/goromdb/storage"
)

type Storage struct {
	db *bdb.BerkeleyDB
}

func New() *Storage {
	return &Storage{nil}
}

func (s *Storage) Load(file string, mux *sync.RWMutex) error {
	db, err := openBDB(file)
	if err != nil {
		return err
	}

	mux.Lock()
	defer mux.Unlock()

	if s.db != nil {
		s.db.Close(0)
	}
	s.db = db
	return nil
}

func (s *Storage) LoadAndIterate(file string, fn storage.IterationFunc, mux *sync.RWMutex) error {
	db, err := openBDB(file)
	if err != nil {
		return err
	}
	err = iterate(db, fn)
	if err != nil {
		return err
	}

	mux.Lock()
	defer mux.Unlock()

	if s.db != nil {
		s.db.Close(0)
	}
	s.db = db
	return nil
}

func openBDB(file string) (*bdb.BerkeleyDB, error) {
	return bdb.OpenBDB(bdb.NoEnv, bdb.NoTxn, file, nil, bdb.BTree, bdb.DbReadOnly, 0)
}

func (s *Storage) Get(key []byte) ([]byte, error) {
	if s.db != nil {
		v, err := s.db.Get(bdb.NoTxn, key, 0)
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
