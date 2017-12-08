package jsonstorage

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/yowcow/goromdb/storage"
)

type Data map[string]string

type Storage struct {
	gzipped bool
	data    Data
}

func New(gzipped bool) *Storage {
	return &Storage{
		gzipped,
		make(Data),
	}
}

func (s *Storage) Load(file string, mux *sync.RWMutex) error {
	data, err := s.openFile(file)
	if err != nil {
		return err
	}

	// Lock, switch, and unlock
	mux.Lock()
	defer mux.Unlock()

	s.data = data
	return nil
}

func (s *Storage) LoadAndIterate(file string, fn storage.IterationFunc, mux *sync.RWMutex) error {
	data, err := s.openFile(file)
	if err != nil {
		return err
	}
	err = iterate(data, fn)
	if err != nil {
		return err
	}

	// Lock, switch, and unlock
	mux.Lock()
	defer mux.Unlock()

	s.data = data
	return nil
}

func (s Storage) openFile(file string) (Data, error) {
	fi, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	r, err := s.newReader(fi)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(r)
	var data Data
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func (s Storage) newReader(rdr io.Reader) (io.Reader, error) {
	if s.gzipped {
		r, err := gzip.NewReader(rdr)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	return rdr, nil
}

func (s Storage) Get(key []byte) ([]byte, error) {
	k := string(key)
	if v, ok := s.data[k]; ok {
		return []byte(v), nil
	}
	return nil, storage.KeyNotFoundError(key)
}

func iterate(data Data, fn storage.IterationFunc) error {
	for k, v := range data {
		if err := fn([]byte(k), []byte(v)); err != nil {
			return err
		}
	}
	return nil
}
