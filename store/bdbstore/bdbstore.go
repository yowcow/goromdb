package bdbstore

import (
	"log"
	"sync"

	"github.com/ajiyoshi-vg/goberkeleydb/bdb"
	"github.com/yowcow/goromdb/store"
)

// Data represents a key-value data
type Data map[string][]byte

// Store represents a store
type Store struct {
	data   Data
	filein <-chan string
	loader *store.Loader
	db     *bdb.BerkeleyDB
	mux    *sync.RWMutex
	logger *log.Logger
}

// New creates a store
func New(filein <-chan string, basedir string, logger *log.Logger) (store.Store, error) {
	data := make(Data)
	loader, err := store.NewLoader(basedir)
	if err != nil {
		return nil, err
	}
	return &Store{
		data,
		filein,
		loader,
		nil,
		new(sync.RWMutex),
		logger,
	}, nil
}

func (s *Store) Start() <-chan bool {
	done := make(chan bool)
	go s.start(done)
	return done
}

func (s *Store) start(done chan<- bool) {
	defer func() {
		s.logger.Println("bdbstore finished")
		close(done)
	}()
	s.logger.Println("bdbstore started")
	for file := range s.filein {
		s.logger.Printf("bdbstore got new file to load at '%s'", file)
		newfile, err := s.loader.DropIn(file)
		if err != nil {
			s.logger.Printf("bdbstore failed dropping file from '%s' into '%s': %s", file, newfile, err.Error())
		} else if err = s.Load(newfile); err != nil {
			s.logger.Printf("bdbstore failed loading data from '%s': %s", newfile, err.Error())
		}
	}
}

func (s *Store) Load(file string) error {
	db, err := openBDB(file)
	if err != nil {
		return err
	}
	s.logger.Printf("bdbstore successfully opened new db at '%s'", file)

	data := make(Data)
	olddb := s.db

	s.mux.Lock()
	s.data = data
	s.db = db
	s.mux.Unlock()

	if olddb != nil {
		olddb.Close(0)
		olddb = nil
		s.logger.Printf("bdbstore successfully closed old db")
	}
	return nil
}

func openBDB(file string) (*bdb.BerkeleyDB, error) {
	return bdb.OpenBDB(bdb.NoEnv, bdb.NoTxn, file, nil, bdb.BTree, bdb.DbReadOnly, 0)
}

// Get retrieves the value for given key from a store
func (s *Store) Get(key []byte) ([]byte, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	k := string(key)
	if v, ok := s.data[k]; ok {
		return v, nil
	} else if s.db != nil {
		v, err := s.db.Get(bdb.NoTxn, key, 0)
		if err != nil {
			s.data[k] = nil
			return nil, err
		}
		s.data[k] = v
		return v, nil
	}
	return nil, store.KeyNotFoundError(key)
}
