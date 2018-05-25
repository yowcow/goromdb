package simplehandler

import (
	"log"

	"github.com/yowcow/goromdb/handler"
	"github.com/yowcow/goromdb/storage"
)

var (
	_ handler.Handler = (*Handler)(nil)
)

// Handler represents a simple handler
type Handler struct {
	StorageHandler
}

// New creates and returns a handler
func New(stg storage.Storage, logger *log.Logger) *Handler {
	return &Handler{StorageHandler{stg, logger}}
}

// Get finds value by given key, and returns the value
func (h *Handler) Get(key []byte) ([]byte, error) {
	return h.storage.Get(key)
}
