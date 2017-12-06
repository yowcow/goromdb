package storage

import (
	"fmt"
	"sync"
)

type Storage interface {
	Get([]byte) ([]byte, error)
	Load(string, *sync.RWMutex) error
	Cursor() (Cursor, error)
}

type Cursor interface {
	Next() ([]byte, []byte, error)
	Close() error
}

func KeyNotFoundError(key []byte) error {
	return fmt.Errorf("key not found error: %s", string(key))
}
