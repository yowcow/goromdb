package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yowcow/go-romdb/protocol"
	"github.com/yowcow/go-romdb/store"
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

func (p TestProtocol) Reply(w *bufio.Writer, key, value []byte) {
	w.Write(key)
	w.WriteRune(' ')
	w.Write(value)
	w.WriteString("\r\n")
}

func (p TestProtocol) Finish(w *bufio.Writer) {
	w.WriteString("BYE\r\n")
}

type TestData map[string]string

type TestStore struct {
	data   TestData
	logger *log.Logger
}

func createTestStore() store.Store {
	data := TestData{
		"foo": "foo!",
		"bar": "bar!!",
	}
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	return &TestStore{data, logger}
}

func (s TestStore) Get(key []byte) ([]byte, error) {
	if v, ok := s.data[string(key)]; ok {
		return []byte(v), nil
	}
	return nil, store.KeyNotFoundError(key)
}

func (s TestStore) Shutdown() error {
	s.logger.Print("store shutting down")
	return nil
}

func TestServer(t *testing.T) {
	type Case struct {
		input    string
		expected []string
	}

	cases := []Case{
		{
			input: "hoge\r\n",
			expected: []string{
				"foo foo!",
				"bar bar!!",
				"BYE",
			},
		},
		{
			input: "fuga\r\n",
			expected: []string{
				"BYE",
			},
		},
	}

	p := createTestProtocol()
	s := createTestStore()
	server := New("tcp", ":11222", p, s)

	go func() {
		server.Start()
	}()

	time.Sleep(1 * time.Second) // should wait server to get started
	conn, err := net.Dial("tcp", "localhost:11222")

	assert.Nil(t, err)

	var output []byte
	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)

	for _, c := range cases {
		w.WriteString(c.input)
		w.Flush()

		for _, expected := range c.expected {
			output, _, err = r.ReadLine()

			assert.Equal(t, expected, string(output))
			assert.Nil(t, err)
		}
	}

	assert.Nil(t, conn.Close())
	assert.Nil(t, server.Shutdown())
}
