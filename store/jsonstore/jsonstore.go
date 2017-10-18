package jsonstore

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

	dataLoaderQuit chan bool
	dataLoaderWg   *sync.WaitGroup
}

func New(file string, logger *log.Logger) store.Store {
	var data Data
	dataUpdate := make(chan Data)

	dataNodeQuit := make(chan bool)
	dataNodeWg := &sync.WaitGroup{}

	dataLoaderQuit := make(chan bool)
	dataLoaderWg := &sync.WaitGroup{}

	s := &Store{
		file,
		data,
		logger,
		dataNodeQuit,
		dataNodeWg,
		dataLoaderQuit,
		dataLoaderWg,
	}

	boot := make(chan bool)

	dataNodeWg.Add(1)
	go s.startDataNode(boot, dataUpdate)
	<-boot

	dataLoaderWg.Add(1)
	go s.startDataLoader(boot, dataUpdate)
	<-boot

	close(boot)

	return s
}

func (s *Store) startDataNode(boot chan<- bool, dataIn <-chan Data) {
	defer s.dataNodeWg.Done()

	if data, err := LoadJSON(s.file); err == nil {
		if err = store.CheckMD5Sum(s.file, s.file+".md5"); err != nil {
			s.logger.Print("-> data node failed checking MD5 sum: ", err)
		} else {
			s.data = data
		}
	} else {
		s.logger.Print("-> data node failed reading data from file: ", err)
	}

	boot <- true
	s.logger.Print("-> data node started!")

	for {
		select {
		case data := <-dataIn:
			s.logger.Print("-> data node updated!")
			s.data = data
		case <-s.dataNodeQuit:
			s.logger.Print("-> data node finished!")
			return
		}
	}
}

func (s Store) startDataLoader(boot chan<- bool, dataOut chan<- Data) {
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
			if data, err := LoadJSON(s.file); err == nil {
				dataOut <- data
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
	if v, ok := s.data[string(key)]; ok {
		return []byte(v), nil
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
