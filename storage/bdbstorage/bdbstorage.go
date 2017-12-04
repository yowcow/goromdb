package bdbstorage

import (
	"github.com/ajiyoshi-vg/goberkeleydb/bdb"
	"github.com/yowcow/goromdb/storage"
)

type Data map[string][]byte

type Storage struct {
	db   *bdb.BerkeleyDB
	data Data
}

func New() *Storage {
	return &Storage{
		nil,
		make(Data),
	}
}

func (s *Storage) Load(file string) error {
	db, err := openBDB(file)
	if err != nil {
		return err
	}
	oldDB := s.db
	data := make(Data)
	s.db = db
	s.data = data
	if oldDB != nil {
		oldDB.Close(0)
		oldDB = nil
	}
	return nil
}

func openBDB(file string) (*bdb.BerkeleyDB, error) {
	return bdb.OpenBDB(bdb.NoEnv, bdb.NoTxn, file, nil, bdb.BTree, bdb.DbReadOnly, 0)
}

func (s *Storage) Get(key []byte) ([]byte, error) {
	v, ok := s.data[string(key)]
	if ok {
		if v != nil {
			return v, nil
		}
		return nil, storage.KeyNotFoundError(key)
	}
	if s.db != nil {
		v, err := s.db.Get(bdb.NoTxn, key, 0)
		if err != nil {
			s.data[string(key)] = nil
			return nil, storage.KeyNotFoundError(key)
		}
		s.data[string(key)] = v
		return v, nil
	}
	return nil, storage.KeyNotFoundError(key)
}
