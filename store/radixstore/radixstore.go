package radixstore

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/armon/go-radix"
	"github.com/yowcow/goromdb/reader"
	"github.com/yowcow/goromdb/store"
)

type Store struct {
	tree             *radix.Tree
	filein           <-chan string
	gzipped          bool
	loader           *store.Loader
	createReaderFunc reader.NewReaderFunc
	mux              *sync.RWMutex
	logger           *log.Logger
}

func New(
	filein <-chan string,
	gzipped bool,
	basedir string,
	createReaderFunc reader.NewReaderFunc,
	logger *log.Logger,
) (store.Store, error) {
	loader, err := store.NewLoader(basedir, "data.csv")
	if err != nil {
		return nil, err
	}
	if createReaderFunc == nil {
		return nil, fmt.Errorf("createReaderFunc cannot be a nil")
	}
	return &Store{
		radix.New(),
		filein,
		gzipped,
		loader,
		createReaderFunc,
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

	ior, err := store.NewReader(f, s.gzipped)
	if err != nil {
		return err
	}

	tree := radix.New()
	r := s.createReaderFunc(ior)
	if err = buildTree(tree, r); err != nil {
		return err
	}

	s.mux.Lock()
	s.tree = tree
	s.mux.Unlock()
	s.logger.Println("radixstore successfully replaced a tree")
	return nil
}

func buildTree(tree *radix.Tree, r reader.Reader) error {
	for {
		k, v, err := r.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		tree.Insert(string(k), v)
	}
}

func (s Store) Get(k []byte) ([]byte, []byte, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()
	prefix, v, ok := s.tree.LongestPrefix(string(k))
	if !ok {
		return nil, nil, store.KeyNotFoundError(k)
	}
	return []byte(prefix), v.([]byte), nil
}
