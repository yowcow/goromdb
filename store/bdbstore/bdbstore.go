package bdbstore

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/ajiyoshi-vg/goberkeleydb/bdb"
	"github.com/yowcow/go-romdb/store"
)

type Data map[string]string

type Store struct {
	file   string
	data   Data
	db     *bdb.BerkeleyDB
	logger *log.Logger

	dataNodeQuit chan bool
	dataNodeWg   *sync.WaitGroup

	dataLoaderQuit chan bool
	dataLoaderWg   *sync.WaitGroup
}

func New(file string) (store.Store, error) {
	var data Data
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	dbUpdate := make(chan *bdb.BerkeleyDB)

	dataNodeQuit := make(chan bool)
	dataNodeWg := &sync.WaitGroup{}

	dataLoaderQuit := make(chan bool)
	dataLoaderWg := &sync.WaitGroup{}

	s := &Store{file, data, nil, logger, dataNodeQuit, dataNodeWg, dataLoaderQuit, dataLoaderWg}

	boot := make(chan bool)

	dataNodeWg.Add(1)
	go s.startDataNode(boot, dbUpdate)
	<-boot

	dataLoaderWg.Add(1)
	go s.startDataLoader(boot, dbUpdate)
	<-boot

	close(boot)

	return s, nil
}

func (s *Store) startDataNode(boot chan<- bool, dbIn <-chan *bdb.BerkeleyDB) {
	defer s.dataNodeWg.Done()

	if db, err := OpenBDB(s.file); err == nil {
		s.db = db
	}

	boot <- true
	s.logger.Print("-> datanode started!")

	for {
		select {
		case newDB := <-dbIn:
			oldDB := s.db
			s.db = newDB
			if oldDB != nil {
				oldDB.Close(0)
				oldDB = nil
			}
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
	update := make(chan bool)
	quit := make(chan bool)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go store.NewWatcher(d, s.file, s.logger, update, quit, wg)

	boot <- true
	s.logger.Print("-> data loader started!")

	for {
		select {
		case <-update:
			if db, err := OpenBDB(s.file); err == nil {
				dbOut <- db
			} else {
				s.logger.Print("-> data loader failed reading data from file: ", err)
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

func (s Store) Get(key []byte) ([]byte, error) {
	if s.db != nil {
		if v, err := s.db.Get(bdb.NoTxn, key, 0); err == nil {
			return v, nil
		} else {
			return nil, err
		}
	}
	return nil, store.KeyNotFoundError(key)
}

func (s Store) Shutdown() error {
	s.dataLoaderQuit <- true
	close(s.dataLoaderQuit)
	s.dataLoaderWg.Wait()

	s.dataNodeQuit <- true
	close(s.dataNodeQuit)
	s.dataNodeWg.Wait()

	return nil
}

func OpenBDB(file string) (*bdb.BerkeleyDB, error) {
	return bdb.OpenBDB(bdb.NoEnv, bdb.NoTxn, file, nil, bdb.BTree, bdb.DbReadOnly, 0)
}
