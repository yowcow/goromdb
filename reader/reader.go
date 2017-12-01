package reader

import (
	"io"
)

// NewReaderFunc is a new reader creating function
type NewReaderFunc func(io.Reader) Reader

// Reader is an interface for general reader
type Reader interface {
	Read() ([]byte, []byte, error)
}
