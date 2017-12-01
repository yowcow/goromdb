package reader

import (
	"io"
)

type NewReaderFunc func(io.Reader) Reader

type Reader interface {
	Read() ([]byte, []byte, error)
}
