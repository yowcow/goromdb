package bdbstore

import (
	"bytes"
	"log"

	"github.com/yowcow/go-romdb/store"
	bdb "github.com/yowcow/go-romdb/store/bdbstore"
	"github.com/yowcow/go-romdb/store/memcachedb"
)

type Store struct {
	proxy  store.Store
	logger *log.Logger
}

func New(file string, logger *log.Logger) store.Store {
	proxy := bdb.New(file, logger)
	s := &Store{proxy, logger}
	return s
}

func (s Store) Get(key []byte) ([]byte, error) {
	val, err := s.proxy.Get(key)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(val)
	_, v, _, err := memcachedb.Deserialize(r)
	if err != nil {
		s.logger.Print(err)
		return nil, err
	}

	return v, nil
}

func (s Store) Shutdown() error {
	return s.proxy.Shutdown()
}
