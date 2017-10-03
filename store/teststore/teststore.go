package teststore

import (
	"github.com/yowcow/go-romdb/store"
)

type Data map[string]string

type StoreTest struct {
	data Data
}

func New() (store.Store, error) {
	data := Data{
		"foo": "my test foo",
		"bar": "my test bar!!",
	}
	return &StoreTest{data}, nil
}

func (s StoreTest) Get(key string) (string, error) {
	if v, ok := s.data[key]; ok {
		return v, nil
	}
	return "", store.KeyNotFoundError(key)
}
