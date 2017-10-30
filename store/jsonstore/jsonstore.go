package jsonstore

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/yowcow/goromdb/store"
)

// Data represents a key-value data
type Data map[string]string

// Store represents a store
type Store struct {
	file   string
	data   Data
	logger *log.Logger
	loader *store.Loader
	quit   chan bool
	wg     *sync.WaitGroup
}

// New creates a new store
func New(file string, logger *log.Logger) store.Store {
	var data Data
	dataUpdate := make(chan Data)

	quit := make(chan bool)
	wg := &sync.WaitGroup{}

	baseDir := filepath.Dir(file)
	loader := store.NewLoader(baseDir, logger)

	if err := loader.BuildStoreDirs(); err != nil {
		logger.Print("-> store failed creating directories: ", err)
	}

	s := &Store{
		file,
		data,
		logger,
		loader,
		quit,
		wg,
	}

	boot := make(chan bool)

	wg.Add(1)
	go s.startDataNode(boot, dataUpdate)

	<-boot
	close(boot)

	return s
}

func (s *Store) startDataNode(boot chan<- bool, dataIn <-chan Data) {
	defer s.wg.Done()

	d := 5 * time.Second
	watcher := store.NewWatcher(s.file, d, store.CheckMD5Sum, s.logger)

	if watcher.IsLoadable() {
		if data, err := LoadData(s.loader, s.file); err != nil {
			s.logger.Print("-> data loader failed: ", err)
		} else {
			s.data = data
		}
	}

	boot <- true
	s.logger.Print("-> data node started!")

	update := make(chan bool)
	quit := make(chan bool)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go watcher.Start(update, quit, wg)

	for {
		select {
		case <-update:
			s.logger.Print("-> data node ready to update!")
			if data, err := LoadData(s.loader, s.file); err != nil {
				s.logger.Print("-> data node failed loading data: ", err)
			} else {
				s.data = data
				s.logger.Print("-> data node succeeded loading new data")
				if err := s.loader.CleanOldDirs(); err != nil {
					s.logger.Print("-> data node failed cleaning old directory: ", err)
				}
				s.logger.Print("-> data node updated!")
			}
		case <-s.quit:
			quit <- true
			close(quit)
			wg.Wait()
			close(update)
			s.logger.Print("-> data node finished!")
			return
		}
	}
}

// Get retrieves the value for given key from a store
func (s Store) Get(key []byte) ([]byte, error) {
	if v, ok := s.data[string(key)]; ok {
		return []byte(v), nil
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

// LoadData moves file into next store dir, unmarshals JSON and returns Data
func LoadData(loader *store.Loader, file string) (Data, error) {
	nextFile, err := loader.MoveFileToNextDir(file)
	if err != nil {
		return nil, err
	}
	data, err := LoadJSON(nextFile)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// LoadJSON reads file and parse JSON into Data
func LoadJSON(file string) (Data, error) {
	var data Data
	fi, err := os.Open(file)
	if err != nil {
		return data, err
	}
	defer fi.Close()
	dec := json.NewDecoder(fi)
	err = dec.Decode(&data)
	if err != nil {
		return data, err
	}
	return data, nil
}
