package bdbstore

import (
	"log"
	"sync"
	"time"

	"github.com/ajiyoshi-vg/goberkeleydb/bdb"
	"github.com/yowcow/go-romdb/store"
)

type Data map[string][]byte

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

func New(file string, logger *log.Logger) store.Store {
	data := make(Data)
	dbUpdate := make(chan *bdb.BerkeleyDB)

	dataNodeQuit := make(chan bool)
	dataNodeWg := &sync.WaitGroup{}

	dataLoaderQuit := make(chan bool)
	dataLoaderWg := &sync.WaitGroup{}

	s := &Store{
		file:           file,
		data:           data,
		db:             nil,
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

	if db, err := OpenBDB(s.file); err == nil {
		if err = store.CheckMD5Sum(s.file, s.file+".md5"); err != nil {
			s.logger.Print("-> data node failed checking MD5 sum: ", err)
		} else {
			s.db = db
		}
	} else {
		s.logger.Print("-> data node failed reading data from file: ", err)
	}

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
	update := make(chan bool)
	quit := make(chan bool)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go store.NewWatcher(d, s.file, s.logger, update, quit, wg, store.CheckMD5Sum)

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
