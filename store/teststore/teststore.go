package teststore

import (
	"github.com/yowcow/go-romdb/store"
)

// Data represents a key-value data
type Data map[string][]byte

// Store represents a data store
type Store struct {
	data Data
}

// New creates a new store
func New() store.Store {
	data := Data{
		"foo": []byte("my test foo"),
		"bar": []byte("my test bar!!"),
	}
	return &Store{data}
}

// Get retrieves the value for given key
func (s Store) Get(key []byte) ([]byte, error) {
	if v, ok := s.data[string(key)]; ok {
		return v, nil
	}
	return nil, store.KeyNotFoundError(key)
}

// Shutdown terminates store
func (s Store) Shutdown() error {
	return nil
}
