package memcdstorage

import (
	"github.com/yowcow/goromdb/storage"
)

var (
	_ storage.NSStorage = (*NSStorage)(nil)
)

type NSStorage struct {
	proxy storage.NSStorage
}

func NewNS(proxy storage.NSStorage) *NSStorage {
	return &NSStorage{proxy}
}

// Load loads data into storage
func (s NSStorage) Load(file string) error {
	return s.proxy.Load(file)
}

// Get finds a given key in storage, deserialize its value into memcachedb format, and returns
func (s NSStorage) Get(key []byte) ([]byte, error) {
	val, err := s.proxy.Get(key)
	if err != nil {
		return nil, err
	}
	return unmarshalMemcachedbBytes(key, val)
}

// Get finds a given ns+key in storage, deserialize its value into memcachedb format, and returns
func (s NSStorage) GetNS(ns, key []byte) ([]byte, error) {
	val, err := s.proxy.GetNS(ns, key)
	if err != nil {
		return nil, err
	}
	return unmarshalMemcachedbBytes(key, val)
}
