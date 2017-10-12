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

	watcherQuit chan bool
	watcherWg   *sync.WaitGroup
}

func New(file string) (store.Store, error) {
	var data Data
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	dbUpdate := make(chan *bdb.BerkeleyDB)

	dataNodeQuit := make(chan bool)
	dataNodeWg := &sync.WaitGroup{}

	watcherQuit := make(chan bool)
	watcherWg := &sync.WaitGroup{}

	s := &Store{file, data, nil, logger, dataNodeQuit, dataNodeWg, watcherQuit, watcherWg}

	boot := make(chan bool)

	dataNodeWg.Add(1)
	go s.startDataNode(boot, dbUpdate)
	<-boot

	watcherWg.Add(1)
	go s.startWatcher(boot, dbUpdate)
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
			oldDB.Close(0)
			s.logger.Print("-> datanode updated!")
		case <-s.dataNodeQuit:
			if s.db != nil {
				s.db.Close(0)
				s.db = nil
			}
			s.logger.Print("-> datanode finished!")
			return
		}
	}
}

func (s Store) startWatcher(boot chan<- bool, dbOut chan<- *bdb.BerkeleyDB) {
	defer s.watcherWg.Done()

	var lastModified time.Time

	if fi, err := os.Stat(s.file); err == nil {
		lastModified = fi.ModTime()
	}

	d := 5 * time.Second
	t := time.NewTimer(d)

	boot <- true
	s.logger.Print("-> watcher started!")

	for {
		select {
		case <-t.C:
			if fi, err := os.Stat(s.file); err == nil {
				if fi.ModTime() != lastModified {
					lastModified = fi.ModTime()
					if db, err := OpenBDB(s.file); err == nil {
						dbOut <- db
					} else {
						s.logger.Print("-> watcher failed reading data from file: ", err)
					}
				}
			} else {
				s.logger.Print("-> watcher file check failed: ", err)
			}
			t.Reset(d)
		case <-s.watcherQuit:
			if !t.Stop() {
				<-t.C
			}
			s.logger.Print("-> watcher finished!")
			return
		}
	}
}

func (s Store) Get(key string) (string, error) {
	if s.db != nil {
		if v, err := s.db.Get(bdb.NoTxn, []byte(key), 0); err == nil {
			return string(v), nil
		} else {
			return "", err
		}
	}
	return "", store.KeyNotFoundError(key)
}

func (s Store) Shutdown() error {
	s.watcherQuit <- true
	close(s.watcherQuit)
	s.watcherWg.Wait()

	s.dataNodeQuit <- true
	close(s.dataNodeQuit)
	s.dataNodeWg.Wait()

	return nil
}

func OpenBDB(file string) (*bdb.BerkeleyDB, error) {
	return bdb.OpenBDB(bdb.NoEnv, bdb.NoTxn, file, nil, bdb.BTree, bdb.DbReadOnly, 0)
}
