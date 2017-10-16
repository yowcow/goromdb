package bdbstore

import (
	"bytes"

	"github.com/yowcow/go-romdb/store"
	bdb "github.com/yowcow/go-romdb/store/bdbstore"
	"github.com/yowcow/go-romdb/store/memcachedb"
)

type Store struct {
	proxy store.Store
}

func New(file string) (store.Store, error) {
	proxy, _ := bdb.New(file)
	s := &Store{proxy}
	return s, nil
}

func (s Store) Get(key []byte) ([]byte, error) {
	val, err := s.proxy.Get(key)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(val)
	_, v, _, err := memcachedb.Deserialize(r)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (s Store) Shutdown() error {
	return s.proxy.Shutdown()
}
