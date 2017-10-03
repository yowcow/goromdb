package jsonstore

import (
	"encoding/json"
	"io/ioutil"

	"github.com/yowcow/go-romdb/store"
)

type Data map[string]string

type JSONStore struct {
	data map[string]string
}

func New(file string) (store.Store, error) {
	var data Data
	b, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &data)

	if err != nil {
		return nil, err
	}

	return &JSONStore{data}, nil
}

func (s JSONStore) Get(key string) (string, error) {
	if v, ok := s.data[key]; ok {
		return v, nil
	}
	return "", store.KeyNotFoundError(key)
}
