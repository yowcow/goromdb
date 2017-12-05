package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goromdb/handler/simplehandler"
	"github.com/yowcow/goromdb/protocol"
	"github.com/yowcow/goromdb/storage"
	"github.com/yowcow/goromdb/testutil"
)

type TestKeywords map[string][][]byte

type TestProtocol struct {
	keywords TestKeywords
}

func createTestProtocol() protocol.Protocol {
	keywords := TestKeywords{
		"hoge": {[]byte("foo"), []byte("bar")},
	}
	return &TestProtocol{keywords}
}

func (p TestProtocol) Parse(line []byte) ([][]byte, error) {
	if words, ok := p.keywords[string(line)]; ok {
		return words, nil
	}
	return [][]byte{}, fmt.Errorf("invalid command")
}

func (p TestProtocol) Reply(w io.Writer, key, value []byte) {
	w.Write(key)
	w.Write([]byte(" "))
	w.Write(value)
	w.Write([]byte("\r\n"))
}

func (p TestProtocol) Finish(w io.Writer) {
	w.Write([]byte("BYE\r\n"))
}

type TestData map[string]string

type TestStorage struct {
	data   TestData
	logger *log.Logger
}

func createTestStorage(logger *log.Logger) storage.Storage {
	data := TestData{
		"foo": "foo!",
		"bar": "bar!!",
	}
	return &TestStorage{data, logger}
}

func (s TestStorage) Load(file string) error {
	return nil
}

func (s TestStorage) Get(key []byte) ([]byte, error) {
	if v, ok := s.data[string(key)]; ok {
		return []byte(v), nil
	}
	return nil, storage.KeyNotFoundError(key)
}

func (s TestStorage) Cursor() (storage.Cursor, error) {
	return nil, fmt.Errorf("cursor not supported")
}

func TestHandleConn(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	p := createTestProtocol()
	stg := createTestStorage(logger)
	h := simplehandler.New(stg, logger)

	sock := filepath.Join(dir, "test.sock")
	svr := New("unix", sock, p, h, logger)

	done := make(chan bool)
	ln, err := net.Listen("unix", sock)
	if err != nil {
		panic(err)
	}
	go func(d chan<- bool) {
		defer func() {
			close(d)
		}()
		for {
			conn, err := ln.Accept()
			if err != nil {
				break
			}
			svr.HandleConn(conn)
		}
	}(done)

	type Case struct {
		input    string
		expected []string
		subtest  string
	}
	cases := []Case{
		{
			"hoge\r\n",
			[]string{
				"foo foo!",
				"bar bar!!",
				"BYE",
			},
			"hoge returns 3 lines of message",
		},
		{
			"fuga\r\n",
			[]string{
				"BYE",
			},
			"fuga returns 1 line of message",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			conn, err := net.Dial("unix", sock)
			if err != nil {
				panic(err)
			}
			defer conn.Close()

			r := bufio.NewReader(conn)
			_, err = conn.Write([]byte(c.input))

			assert.Nil(t, err)

			for _, row := range c.expected {
				actual, _, err := r.ReadLine()

				assert.Nil(t, err)
				assert.Equal(t, row, string(actual))
			}
		})
	}

	ln.Close()
	<-done
}
