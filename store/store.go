package store

import (
	"compress/gzip"
	"fmt"
	"io"
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

// NewReader returns a gzip reader if gzipped is true
func NewReader(f io.Reader, gzipped bool) (io.Reader, error) {
	if gzipped {
		return gzip.NewReader(f)
	}
	return f, nil
}
