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
	logger *log.Logger
	loader *store.Loader

	dataNodeQuit chan bool
	dataNodeWg   *sync.WaitGroup

	dataLoaderQuit chan bool
	dataLoaderWg   *sync.WaitGroup
}

// New creates a store
func New(file string, logger *log.Logger) store.Store {
	data := make(Data)
	dbUpdate := make(chan *bdb.BerkeleyDB)

	dataNodeQuit := make(chan bool)
	dataNodeWg := &sync.WaitGroup{}

	dataLoaderQuit := make(chan bool)
	dataLoaderWg := &sync.WaitGroup{}

	baseDir := filepath.Dir(file)
	loader := store.NewLoader(baseDir, logger)

	if err := loader.BuildStoreDirs(); err != nil {
		logger.Print("-> store failed creating directories: ", err)
	}

	s := &Store{
		file:           file,
		data:           data,
		db:             nil,
		loader:         loader,
		logger:         logger,
		dataNodeQuit:   dataNodeQuit,
		dataNodeWg:     dataNodeWg,
		dataLoaderQuit: dataLoaderQuit,
		dataLoaderWg:   dataLoaderWg,
	}

	boot := make(chan bool)

	dataNodeWg.Add(1)
	go s.startDataNode(boot, dbUpdate)
	<-boot

	dataLoaderWg.Add(1)
	go s.startDataLoader(boot, dbUpdate)
	<-boot

	close(boot)

	return s
}

func (s *Store) startDataNode(boot chan<- bool, dbIn <-chan *bdb.BerkeleyDB) {
	defer s.dataNodeWg.Done()

	boot <- true
	s.logger.Print("-> data node started!")

	for {
		select {
		case newDB := <-dbIn:
			oldDB := s.db
			s.db = newDB
			if oldDB != nil {
				oldDB.Close(0)
				oldDB = nil
				if err := s.loader.CleanOldDirs(); err != nil {
					s.logger.Print("-> data node failed cleaning old directory: ", err)
				}
			}
			s.data = make(Data)
			s.logger.Print("-> data node updated!")
		case <-s.dataNodeQuit:
			if s.db != nil {
				s.db.Close(0)
				s.db = nil
			}
			s.logger.Print("-> data node finished!")
			return
		}
	}
}

func (s Store) startDataLoader(boot chan<- bool, dbOut chan<- *bdb.BerkeleyDB) {
	defer s.dataLoaderWg.Done()

	d := 5 * time.Second
	watcher := store.NewWatcher(s.file, d, store.CheckMD5Sum, s.logger)

	if watcher.IsLoadable() {
		if nextFile, err := s.loader.MoveFileToNextDir(s.file); err != nil {
			s.logger.Print("-> data loader failed moving file to store directory: ", err)
		} else {
			s.logger.Print("-> data loader reading data from file: ", nextFile)
			if db, err := OpenBDB(nextFile); err != nil {
				s.logger.Print("-> data loader failed reading data from file: ", err)
			} else {
				dbOut <- db
			}
		}
	}

	boot <- true
	s.logger.Print("-> data loader started!")

	update := make(chan bool)
	quit := make(chan bool)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go watcher.Start(update, quit, wg)

	for {
		select {
		case <-update:
			if nextFile, err := s.loader.MoveFileToNextDir(s.file); err != nil {
				s.logger.Print("-> data loader failed moving file to store directory: ", err)
			} else {
				s.logger.Print("-> data loader reading data from file: ", nextFile)
				if db, err := OpenBDB(nextFile); err != nil {
					s.logger.Print("-> data loader failed reading data from file: ", err)
				} else {
					dbOut <- db
				}
			}
		case <-s.dataLoaderQuit:
			quit <- true
			close(quit)
			wg.Wait()
			close(update)
			s.logger.Print("-> data loader finished!")
			return
		}
	}
}

// Get retrieves the value for given key from a store
func (s Store) Get(key []byte) ([]byte, error) {
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
	s.dataLoaderQuit <- true
	close(s.dataLoaderQuit)
	s.dataLoaderWg.Wait()

	s.dataNodeQuit <- true
	close(s.dataNodeQuit)
	s.dataNodeWg.Wait()

	return nil
}

// OpenBDB creates a BDB handle for given file
func OpenBDB(file string) (*bdb.BerkeleyDB, error) {
	return bdb.OpenBDB(bdb.NoEnv, bdb.NoTxn, file, nil, bdb.BTree, bdb.DbReadOnly, 0)
}
