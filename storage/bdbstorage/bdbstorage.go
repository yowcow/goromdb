package bdbstorage

import (
	"fmt"
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

	// Lock, switch, and unlock
	mux.Lock()
	defer mux.Unlock()
	oldDB := s.db
	s.db = db
	if oldDB != nil {
		if err = oldDB.Close(0); err != nil {
			return err
		}
		oldDB = nil
	}
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

func (s Storage) Iterate(fn storage.IteratorFunc) error {
	if s.db == nil {
		return fmt.Errorf("no bdb handle in storage")
	}

	cur, err := s.db.NewCursor(bdb.NoTxn, 0)
	if err != nil {
		return nil
	}

	for k, v, err := cur.First(); err == nil; k, v, err = cur.Next() {
		if err = fn(k, v); err != nil {
			return err
		}
	}

	return cur.Close()
}
