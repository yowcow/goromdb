package bdbstorage

import (
	"github.com/yowcow/goromdb/storage"
)

var (
	_ storage.NSStorage = (*NSStorage)(nil)
)

// NSStorage represents a namespaced BDB storage
type NSStorage struct {
	Storage
}

// NewNS creates and returns a storage
func NewNS() *NSStorage {
	return &NSStorage{*New()}
}

// GetNS finds a given ns+key in db, and returns its value
func (s *NSStorage) GetNS(ns, key []byte) ([]byte, error) {
	fullKey := make([]byte, 0, len(ns)+len(key))
	fullKey = append(fullKey, ns...)
	fullKey = append(fullKey, key...)
	return s.Get(fullKey)
}
