package storage

import (
	"fmt"
)

// IterationFunc defines an interface to a callback function for LoadAndIterate()
type IterationFunc func([]byte, []byte) error

// Storage defines an interface to a storage
type Storage interface {
	Get([]byte) ([]byte, error)
	Load(string) error
	LoadAndIterate(string, IterationFunc) error
}

// KeyNotFoundError returns an error for key-not-found
func KeyNotFoundError(key []byte) error {
	return fmt.Errorf("key not found error: %s", string(key))
}
