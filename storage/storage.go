package storage

import (
	"fmt"
	"sync"
)

type IterationFunc func([]byte, []byte) error

type Storage interface {
	Get([]byte) ([]byte, error)
	Load(string, *sync.RWMutex) error
	LoadAndIterate(string, IterationFunc, *sync.RWMutex) error
}

func KeyNotFoundError(key []byte) error {
	return fmt.Errorf("key not found error: %s", string(key))
}
