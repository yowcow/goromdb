package boltstorage

import (
	"sync"
	"sync/atomic"

	"github.com/yowcow/goromdb/storage"
)

var (
	_ storage.NSStorage = (*NSStorage)(nil)
)

// NSStorage represents a namespaced BoldDB storage
type NSStorage struct {
	Storage
}

// NewNS creates and returns a storage
func NewNS() *NSStorage {
	return &NSStorage{Storage{new(atomic.Value), nil, new(sync.RWMutex)}}
}

// GetNS finds a given bucket and key in db, and returns its value
func (s *NSStorage) GetNS(ns, key []byte) ([]byte, error) {
	if ns == nil {
		return nil, storage.InternalError("please specify bucket")
	}

	s.mux.RLock()
	defer s.mux.RUnlock()

	db := s.getDB()
	if db == nil {
		return nil, storage.InternalError("couldn't load db")
	}

	return getFromBucket(db, ns, key)
}
