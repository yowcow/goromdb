package handler

import (
	"github.com/yowcow/goromdb/loader"
)

// Handler defines an interface to a handler
type Handler interface {
	Start(<-chan string, *loader.Loader) <-chan bool
	Load(string) error
	Get(key []byte) ([]byte, error)
}

// NSHandler defines an interface to a handler with namespace support
type NSHandler interface {
	Handler
	GetNS(ns, key []byte) ([]byte, error)
}
