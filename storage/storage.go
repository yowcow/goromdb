package storage

import (
	"fmt"
)

type Storage interface {
	Get([]byte) ([]byte, error)
	Load(string) error
	Cursor() (Cursor, error)
}

type Cursor interface {
	Next() ([]byte, []byte, error)
	Close() error
}

func KeyNotFoundError(key []byte) error {
	return fmt.Errorf("key not found error: %s", string(key))
}
