package bdbstore

import (
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/ajiyoshi-vg/goberkeleydb/bdb"
	"github.com/yowcow/goromdb/store"
)

// Data represents a key-value data
type Data map[string][]byte

// Store represents a store
type Store struct {
	file   string
	data   Data
	db     *bdb.BerkeleyDB
	mx     *sync.RWMutex
	logger *log.Logger
	loader *store.Loader
	quit   chan bool
	wg     *sync.WaitGroup
}

// New creates a store
func New(file string, logger *log.Logger) store.Store {
	data := make(Data)
	mx := new(sync.RWMutex)
	dbUpdate := make(chan *bdb.BerkeleyDB)

	quit := make(chan bool)
	wg := &sync.WaitGroup{}

	baseDir := filepath.Dir(file)
	loader := store.NewLoader(baseDir, logger)

	if err := loader.BuildStoreDirs(); err != nil {
		logger.Print("store failed creating directories: ", err)
	}

	s := &Store{
		file:   file,
		data:   data,
		mx:     mx,
		db:     nil,
		loader: loader,
		logger: logger,
		quit:   quit,
		wg:     wg,
	}

	boot := make(chan bool)

	wg.Add(1)
	go s.startDataNode(boot, dbUpdate)

	<-boot
	close(boot)

	return s
}

func (s *Store) startDataNode(boot chan<- bool, dbIn <-chan *bdb.BerkeleyDB) {
	defer s.wg.Done()

	d := 5 * time.Second
	watcher := store.NewWatcher(s.file, d, store.CheckMD5Sum, s.logger)

	if watcher.IsLoadable() {
		if db, err := LoadBDB(s.loader, s.file); err != nil {
			s.logger.Print("data node failed loading a new db: ", err)
		} else {
			s.db = db
		}
	}

	boot <- true
	s.logger.Print("data node started")

	update := make(chan bool)
	quit := make(chan bool)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go watcher.Start(update, quit, wg)

	for {
		select {
		case <-update:
			s.logger.Print("data node found update")
			if db, err := LoadBDB(s.loader, s.file); err != nil {
				s.logger.Print("data node failed loading a new db: ", err)
			} else {
				s.mx.Lock()
				oldDB := s.db
				s.db = db
				s.data = make(Data)
				s.mx.Unlock()
				s.logger.Print("data node succeeded loading a new db")
				if oldDB != nil {
					oldDB.Close(0)
					oldDB = nil
					s.logger.Print("data node succeeded closing an old db")
					if err := s.loader.CleanOldDirs(); err != nil {
						s.logger.Print("data node failed cleaning old directories: ", err)
					}
				}
				s.logger.Print("data node updated")
			}
		case <-s.quit:
			if s.db != nil {
				s.db.Close(0)
				s.db = nil
			}
			s.logger.Print("data node finished")
			return
		}
	}
}

// Get retrieves the value for given key from a store
func (s *Store) Get(key []byte) ([]byte, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()
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

// Shutdown terminates a store
func (s Store) Shutdown() error {
	s.quit <- true
	close(s.quit)
	s.wg.Wait()
	return nil
}

// LoadBDB moves file into next store dir, opens BDB handle, and returns
func LoadBDB(loader *store.Loader, file string) (*bdb.BerkeleyDB, error) {
	nextFile, err := loader.MoveFileToNextDir(file)
	if err != nil {
		return nil, err
	}
	db, err := OpenBDB(nextFile)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// OpenBDB creates a BDB handle for given file
func OpenBDB(file string) (*bdb.BerkeleyDB, error) {
	return bdb.OpenBDB(bdb.NoEnv, bdb.NoTxn, file, nil, bdb.BTree, bdb.DbReadOnly, 0)
}
