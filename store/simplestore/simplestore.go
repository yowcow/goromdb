package simplestore

import (
	"log"
	"sync"

	"github.com/yowcow/goromdb/loader"
	"github.com/yowcow/goromdb/storage"
	"github.com/yowcow/goromdb/store"
)

type Store struct {
	filein  <-chan string
	loader  *loader.Loader
	storage storage.Storage
	logger  *log.Logger
	mux     *sync.RWMutex
}

func New(filein <-chan string, ldr *loader.Loader, stg storage.Storage, logger *log.Logger) store.Store {
	return &Store{
		filein,
		ldr,
		stg,
		logger,
		new(sync.RWMutex),
	}
}

func (s Store) Start() <-chan bool {
	done := make(chan bool)
	go s.start(done)
	return done
}

func (s Store) start(done chan<- bool) {
	defer func() {
		s.logger.Println("simplestore finished")
		close(done)
	}()
	s.logger.Println("simplestore started")
	if newfile, ok := s.loader.FindAny(); ok {
		if err := s.Load(newfile); err != nil {
			s.logger.Printf("simplestore failed loading data from '%s': %s", newfile, err.Error())
		}
	}
	for file := range s.filein {
		s.logger.Printf("simplestore got a new file to load at '%s'", file)
		newfile, err := s.loader.DropIn(file)
		if err != nil {
			s.logger.Printf("simplestore failed dropping file from '%s' into '%s': %s", file, newfile, err.Error())
			continue
		}
		s.logger.Printf("simplestore loading data from '%s'", newfile)
		err = s.Load(newfile)
		if err != nil {
			s.logger.Printf("simplestore failed loading data from '%s': %s", newfile, err.Error())
			continue
		}
		s.logger.Printf("simplestore successfully loaded data from '%s'", newfile)
		if ok := s.loader.CleanUp(); ok {
			s.logger.Print("simplestore successfully removed previously loaded file")
		}
	}
}

func (s Store) Load(file string) error {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.storage.Load(file)
}

func (s Store) Get(key []byte) ([]byte, []byte, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	val, err := s.storage.Get(key)
	if err != nil {
		return nil, nil, err
	}
	return key, val, nil
}
