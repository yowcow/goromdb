package store

import (
	"fmt"
)

type Store interface {
	Get(string) (string, error)
	Shutdown() error
}

func KeyNotFoundError(key string) error {
	return fmt.Errorf("key '%s' not found", key)
}
