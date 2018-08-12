package simplehandler

import (
	"log"

	"github.com/yowcow/goromdb/loader"
	"github.com/yowcow/goromdb/storage"
)

// StorageHandler represents a wrapper to storage.Storage
type StorageHandler struct {
	storage storage.Storage
	logger  *log.Logger
}

// Start starts a handler goroutine
func (h *StorageHandler) Start(filein <-chan string, l *loader.Loader) <-chan bool {
	done := make(chan bool)
	go h.start(filein, l, done)
	return done
}

func (h *StorageHandler) start(filein <-chan string, l *loader.Loader, done chan<- bool) {
	defer func() {
		h.logger.Println("simplehandler finished")
		close(done)
	}()
	h.logger.Println("simplehandler started")
	if newfile, ok := l.FindAny(); ok {
		if err := h.Load(newfile); err != nil {
			h.logger.Printf("simplehandler failed loading data from '%s': %s", newfile, err.Error())
		} else {
			h.logger.Printf("simplehandler loaded data from '%s'", newfile)
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
func (h *StorageHandler) Load(file string) error {
	return h.storage.Load(file)
}
