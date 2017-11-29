package store

import (
	"fmt"
)

// Store represents an interface for a store
type Store interface {
	Start() <-chan bool
	Load(string) error
	Get([]byte) ([]byte, error)
}

// KeyNotFoundError returns an error for key is not found
func KeyNotFoundError(key []byte) error {
	return fmt.Errorf("key '%s' not found", string(key))
}
