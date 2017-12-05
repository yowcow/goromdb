package storage

import (
	"fmt"
)

type Storage interface {
	Get([]byte) ([]byte, error)
	Load(string) error
	AllKeys() [][]byte
}

func KeyNotFoundError(key []byte) error {
	return fmt.Errorf("key not found error: %s", string(key))
}
