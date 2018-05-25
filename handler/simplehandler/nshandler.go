package simplehandler

import (
	"log"

	"github.com/yowcow/goromdb/handler"
	"github.com/yowcow/goromdb/storage"
)

var (
	_ handler.NSHandler = (*NSHandler)(nil)
)

// NSHandler represents a simple namespaced handler
type NSHandler struct {
	Handler
	nsstorage storage.NSStorage
}

// NewNS create a handler with namespace storage
func NewNS(stg storage.NSStorage, logger *log.Logger) *NSHandler {
	return &NSHandler{Handler{StorageHandler{stg, logger}}, stg}
}

// GetNS finds value in namespace by given key, and returns the value
func (h *NSHandler) GetNS(ns, key []byte) ([]byte, error) {
	return h.nsstorage.GetNS(ns, key)
}
