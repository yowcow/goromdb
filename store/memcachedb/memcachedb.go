package memcachedb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"

	"github.com/yowcow/goromdb/store"
)

const _Zero uint8 = 0

// Store represents a store for memcachedb
type Store struct {
	proxy  store.Store
	logger *log.Logger
}

// New creates a store
func New(proxy store.Store, logger *log.Logger) (store.Store, error) {
	return &Store{proxy, logger}, nil
}

// Start starts backend store goroutine, and returns their error
func (s Store) Start() <-chan bool {
	return s.proxy.Start()
}

// Load loads file into backend store, and returns their error
func (s Store) Load(file string) error {
	return s.proxy.Load(file)
}

// Get returns a value from backend store after deserializing it into usable value
func (s Store) Get(key []byte) ([]byte, error) {
	val, err := s.proxy.Get(key)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(val)
	_, v, _, err := Deserialize(r)
	if err != nil {
		s.logger.Print("failed deserializing a value for key '", string(key), "' with error: ", err)
		return nil, store.KeyNotFoundError(key)
	}
	return v, nil
}

// Serialize serializes given key and value into MemcacheDB format binary, and writes to writer
func Serialize(w io.Writer, key, val []byte) error {
	nKey := len(key)
	nBytes := len(val) + 2

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, " %d %d\r\n", 0, len(val))

	sSuffix := buf.Bytes()
	nSuffix := len(sSuffix)

	var data = []interface{}{
		int32(nBytes),
		uint8(nSuffix),
		uint8(nKey),
		_Zero,
		_Zero,
		key,
		_Zero,
		sSuffix,
		val,
		[]byte("\r\n"),
	}

	for _, v := range data {
		var err error
		if err = binary.Write(w, binary.LittleEndian, v); err != nil {
			return fmt.Errorf("failed writing memcachedb binary: %s", err.Error())
		}
	}

	return nil
}

// Deserialize deserializes MemcacheDB format binary from reader into key, value and value length
func Deserialize(r io.Reader) ([]byte, []byte, int, error) {
	var err error
	var (
		nBytes  int32
		nSuffix uint8
		nKey    uint8
		pad1    uint8
		pad2    uint8
	)
	var headers = []interface{}{
		&nBytes,
		&nSuffix,
		&nKey,
		&pad1,
		&pad2,
	}
	for _, v := range headers {
		err = binary.Read(r, binary.LittleEndian, v)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("failed reading memcachedb binary headers: %s", err.Error())
		}
	}

	var (
		key     = make([]byte, nKey)
		sSuffix = make([]byte, nSuffix)
		val     = make([]byte, nBytes-2)
		pad3    uint8
	)
	var body = []interface{}{
		&key,
		&pad3,
		&sSuffix,
		&val,
	}
	for _, v := range body {
		err = binary.Read(r, binary.LittleEndian, v)
		if err != nil {
			return nil, nil, 0, fmt.Errorf("failed reading memcachedb binary body: %s", err.Error())
		}
	}

	return key, val, int(nBytes - 2), nil
}
