package storage

import (
	"fmt"
)

type IterationFunc func([]byte, []byte) error

type Storage interface {
	Get([]byte) ([]byte, error)
	Load(string) error
	LoadAndIterate(string, IterationFunc) error
}

func KeyNotFoundError(key []byte) error {
	return fmt.Errorf("key not found error: %s", string(key))
}
