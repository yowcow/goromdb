package teststore

import (
	"github.com/yowcow/go-romdb/store"
)

type Data map[string][]byte

type StoreTest struct {
	data Data
}

func New() (store.Store, error) {
	data := Data{
		"foo": []byte("my test foo"),
		"bar": []byte("my test bar!!"),
	}
	return &StoreTest{data}, nil
}

func (s StoreTest) Get(key []byte) ([]byte, error) {
	if v, ok := s.data[string(key)]; ok {
		return v, nil
	}
	return nil, store.KeyNotFoundError(key)
}

func (s StoreTest) Shutdown() error {
	return nil
}
