package memcdstorage

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/yowcow/goromdb/storage"
)

const _Zero uint8 = 0

type Storage struct {
	proxy storage.Storage
}

func New(proxy storage.Storage) *Storage {
	return &Storage{proxy}
}

func (s Storage) Load(file string, mux *sync.RWMutex) error {
	return s.proxy.Load(file, mux)
}

func (s Storage) Get(key []byte) ([]byte, error) {
	val, err := s.proxy.Get(key)
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(val)
	_, v, _, err := Deserialize(r)
	if err != nil {
		return nil, storage.KeyNotFoundError(key)
	}
	return v, nil
}

func (s Storage) Iterate(fn storage.IteratorFunc) error {
	return s.proxy.Iterate(fn)
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
