package jsonstore

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/yowcow/go-romdb/store"
)

type Data map[string]string

type Store struct {
	file   string
	data   Data
	logger *log.Logger

	dataNodeQuit chan bool
	dataNodeWg   *sync.WaitGroup

	watcherQuit chan bool
	watcherWg   *sync.WaitGroup
}

func New(file string) (store.Store, error) {
	var data Data
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	dataUpdate := make(chan Data)

	dataNodeWg := &sync.WaitGroup{}
	dataNodeQuit := make(chan bool)

	watcherWg := &sync.WaitGroup{}
	watcherQuit := make(chan bool)

	s := &Store{file, data, logger, dataNodeQuit, dataNodeWg, watcherQuit, watcherWg}

	boot := make(chan bool)

	dataNodeWg.Add(1)
	go s.startDataNode(boot, dataUpdate, dataNodeQuit, dataNodeWg)
	<-boot

	watcherWg.Add(1)
	go s.startWatcher(boot, dataUpdate, watcherQuit, watcherWg)
	<-boot

	close(boot)

	return s, nil
}

func (s Store) Get(key string) (string, error) {
	if v, ok := s.data[key]; ok {
		return v, nil
	}
	return "", store.KeyNotFoundError(key)
}

func (s Store) Shutdown() error {
	s.watcherQuit <- true
	s.watcherWg.Wait()
	close(s.watcherQuit)

	s.dataNodeQuit <- true
	s.dataNodeWg.Wait()
	close(s.dataNodeQuit)

	return nil
}

func (s *Store) startDataNode(boot chan bool, in chan Data, q chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	if data, err := LoadJSON(s.file); err == nil {
		s.data = data
	}

	boot <- true
	s.logger.Print("-> datastore started!")

	for {
		select {
		case data := <-in:
			s.logger.Print("-> datastore updated!")
			s.data = data
		case <-q:
			s.logger.Print("-> datastore finished!")
			return
		}
	}
}

func (s Store) startWatcher(boot chan bool, out chan Data, q chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

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
					if data, err := LoadJSON(s.file); err == nil {
						out <- data
					} else {
						s.logger.Print("-> watcher failed reading data from file: ", err)
					}
				}
			} else {
				s.logger.Print("-> watcher file check failed: ", err)
			}
			t.Reset(d)
		case <-q:
			if !t.Stop() {
				<-t.C
			}
			s.logger.Print("-> watcher finished!")
			return
		}
	}
}

func LoadJSON(file string) (Data, error) {
	var data Data

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
