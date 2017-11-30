package radixstore

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/armon/go-radix"
	"github.com/yowcow/goromdb/store"
)

type Store struct {
	tree   *radix.Tree
	filein <-chan string
	loader *store.Loader
	mux    *sync.RWMutex
	logger *log.Logger
}

func New(filein <-chan string, basedir string, logger *log.Logger) (store.Store, error) {
	loader, err := store.NewLoader(basedir, "data.csv")
	if err != nil {
		return nil, err
	}
	return &Store{
		radix.New(),
		filein,
		loader,
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
		s.logger.Print("radixstore finished")
		close(done)
	}()
	s.logger.Print("radixstore started")
	if file, ok := s.loader.FindAny(); ok {
		if err := s.Load(file); err != nil {
			s.logger.Printf("radixstore failed loading data from '%s': %s", file, err.Error())
		}
	}
	for file := range s.filein {
		s.logger.Printf("radixstore got a new file to load at '%s'", file)
		newfile, err := s.loader.DropIn(file)
		if err != nil {
			s.logger.Printf("radixstore failed dropping file from '%s' into '%s': %s", file, newfile, err.Error())
		} else if err = s.Load(newfile); err != nil {
			s.logger.Printf("radixstore failed loading data from '%s': %s", newfile, err.Error())
		} else if ok := s.loader.CleanUp(); ok {
			s.logger.Print("radixstore successfully removed previously loaded file")
		}
	}
}

func (s *Store) Load(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	s.logger.Printf("radixstore successfully opened a new file at '%s'", file)

	tree := radix.New()
	r := csv.NewReader(f)
	if err = buildTree(tree, r); err != nil {
		return err
	}
	s.mux.Lock()
	s.tree = tree
	s.mux.Unlock()
	s.logger.Println("radixstore successfully replaced a tree")
	return nil
}

func buildTree(tree *radix.Tree, r *csv.Reader) error {
	for {
		record, err := r.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		} else if len(record) != 2 {
			return fmt.Errorf("radixstore cannot load a row with a number of elements not exactly 2: %d", len(record))
		}
		tree.Insert(record[0], record[1])
	}
}

func (s Store) Get(key []byte) ([]byte, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	_, v, ok := s.tree.LongestPrefix(string(key))
	if !ok {
		return nil, store.KeyNotFoundError(key)
	}
	return []byte(v.(string)), nil
}
