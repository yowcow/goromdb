package storage

import (
	"fmt"
	"sync"
)

type IteratorFunc func([]byte, []byte) error

type Storage interface {
	Get([]byte) ([]byte, error)
	Load(string, *sync.RWMutex) error
	Iterate(IteratorFunc) error
}

func KeyNotFoundError(key []byte) error {
	return fmt.Errorf("key not found error: %s", string(key))
}
