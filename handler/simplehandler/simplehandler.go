package simplehandler

import (
	"log"

	"github.com/yowcow/goromdb/handler"
	"github.com/yowcow/goromdb/loader"
	"github.com/yowcow/goromdb/storage"
)

type Handler struct {
	storage storage.Storage
	logger  *log.Logger
}

func New(stg storage.Storage, logger *log.Logger) handler.Handler {
	return &Handler{
		stg,
		logger,
	}
}

func (h Handler) Start(filein <-chan string, l *loader.Loader) <-chan bool {
	done := make(chan bool)
	go h.start(filein, l, done)
	return done
}

func (h Handler) start(filein <-chan string, l *loader.Loader, done chan<- bool) {
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

func (h Handler) Load(file string) error {
	return h.storage.Load(file)
}

func (h Handler) Get(key []byte) ([]byte, []byte, error) {
	val, err := h.storage.Get(key)
	if err != nil {
		return nil, nil, err
	}
	return key, val, nil
}
