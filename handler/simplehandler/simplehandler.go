package simplehandler

import (
	"log"

	"github.com/yowcow/goromdb/handler"
	"github.com/yowcow/goromdb/loader"
	"github.com/yowcow/goromdb/storage"
)

var (
	_ handler.Handler   = (*Handler)(nil)
	_ handler.NSHandler = (*Handler)(nil)
)

// Handler represents a simple handler
type Handler struct {
	storage   storage.Storage
	logger    *log.Logger
	nsStorage storage.NSStorage
}

// New creates and returns a handler
func New(stg storage.Storage, logger *log.Logger) *Handler {
	return &Handler{
		storage: stg,
		logger:  logger,
	}
}

// NewWithNS create a handler with namespace storage
func NewWithNS(stg storage.NSStorage, logger *log.Logger) *Handler {
	return &Handler{
		nsStorage: stg,
		logger:    logger,
	}
}

// Start starts a handler goroutine
func (h *Handler) Start(filein <-chan string, l *loader.Loader) <-chan bool {
	done := make(chan bool)
	go h.start(filein, l, done)
	return done
}

func (h *Handler) start(filein <-chan string, l *loader.Loader, done chan<- bool) {
	defer func() {
		h.logger.Println("simplehandler finished")
		close(done)
	}()
	h.logger.Println("simplehandler started")
	if newfile, ok := l.FindAny(); ok {
		if err := h.Load(newfile); err != nil {
			h.logger.Printf("simplehandler failed loading data from '%s': %s", newfile, err.Error())
		}
	}
	for file := range filein {
		h.logger.Printf("simplehandler got a new file to load at '%s'", file)
		newfile, err := l.DropIn(file)
		if err != nil {
			h.logger.Printf("simplehandler failed dropping file from '%s' into '%s': %s", file, newfile, err.Error())
			continue
		}

		h.logger.Printf("simplehandler loading data from '%s'", newfile)
		err = h.Load(newfile)
		if err != nil {
			h.logger.Printf("simplehandler failed loading data from '%s': %s", newfile, err.Error())
			continue
		}

		h.logger.Printf("simplehandler successfully loaded data from '%s'", newfile)
		if ok := l.CleanUp(); ok {
			h.logger.Print("simplehandler successfully removed previously loaded file")
		}
	}
}

// Load loads data into storage
func (h *Handler) Load(file string) error {
	return h.storage.Load(file)
}

// Get finds value by given key, and returns key and value
func (h *Handler) Get(key []byte) ([]byte, []byte, error) {
	val, err := h.storage.Get(key)
	if err != nil {
		return nil, nil, err
	}
	return key, val, nil
}

// GetNS finds value by given key with namespace, and returns key and value
func (h *Handler) GetNS(ns, key []byte) ([]byte, error) {
	if h.nsStorage == nil {
		return nil, storage.InternalError("given storage does not support namespace storage")
	}
	return h.nsStorage.GetNS(ns, key)
}
