package jsonstorage

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"os"

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

func (s *Storage) Load(file string) error {
	fi, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fi.Close()

	r, err := s.newReader(fi)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(r)
	var data Data
	if err := decoder.Decode(&data); err != nil {
		return err
	}
	s.data = data
	return nil
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

func (s Storage) Cursor() (storage.Cursor, error) {
	size := len(s.data)
	if size == 0 {
		return nil, fmt.Errorf("no data in storage")
	}
	keys := make([]string, size)
	i := 0
	for k, _ := range s.data {
		keys[i] = k
		i++
	}
	return &Cursor{s.data, keys, size, 0}, nil
}

type Cursor struct {
	data     Data
	keys     []string
	size     int
	curindex int
}

func (c *Cursor) Next() ([]byte, []byte, error) {
	if c.curindex >= c.size {
		return nil, nil, fmt.Errorf("end of items")
	}
	key := c.keys[c.curindex]
	val := c.data[key]
	c.curindex++
	return []byte(key), []byte(val), nil
}

func (c Cursor) Close() error {
	return nil
}
