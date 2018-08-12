package jsonstorage

import (
	"github.com/yowcow/goromdb/storage"
)

// NSStorage represents a namespaced JSON storage
type NSStorage struct {
	Storage
}

var (
	_ storage.NSStorage = (*NSStorage)(nil)
)

// NewNS creates and returns a namespaced JSON storage
func NewNS(gzipped bool) *NSStorage {
	return &NSStorage{*New(gzipped)}
}

// GetNS finds a given key in given namespace, and returns its value
func (s *NSStorage) GetNS(ns, key []byte) ([]byte, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	ptr := s.data.Load()
	if ptr == nil {
		return nil, storage.KeyNotFoundError(key)
	}

	data := ptr.(Data)
	nsdata, ok := data[string(ns)]
	if !ok {
		return nil, storage.KeyNotFoundError(ns)
	}

	v, ok := nsdata.(map[string]interface{})[string(key)]
	if !ok {
		return nil, storage.KeyNotFoundError(key)
	}

	return []byte(v.(string)), nil
}
