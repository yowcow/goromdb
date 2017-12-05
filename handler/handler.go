package handler

import (
	"github.com/yowcow/goromdb/loader"
)

// Handler defines an interface to a handler
type Handler interface {
	Start(<-chan string, *loader.Loader) <-chan bool
	Load(string) error
	Get([]byte) ([]byte, []byte, error)
}
