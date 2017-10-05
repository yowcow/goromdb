package server

import (
	"bufio"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yowcow/go-romdb/protocol"
	"github.com/yowcow/go-romdb/store"
)

type TestProtocol struct {
}

func createTestProtocol() protocol.Protocol {
	return &TestProtocol{}
}

func (p TestProtocol) Parse(line []byte) ([][]byte, error) {
	if string(line) == "hoge" {
		return [][]byte{[]byte("foo"), []byte("bar")}, nil
	}
	return [][]byte{}, fmt.Errorf("invalid command")
}

func (p TestProtocol) Reply(w *bufio.Writer, key, value string) {
	w.WriteString(key)
	w.WriteRune(' ')
	w.WriteString(value)
	w.WriteString("\r\n")
}

func (p TestProtocol) Finish(w *bufio.Writer) {
	w.WriteString("BYE\r\n")
}

type TestStore struct {
}

func createTestStore() store.Store {
	return &TestStore{}
}

func (s TestStore) Get(key string) (string, error) {
	switch key {
	case "foo":
		return "foo!", nil
	case "bar":
		return "bar!!", nil
	default:
		return "", fmt.Errorf("invalid key")
	}
}

func (s TestStore) Shutdown() error {
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
}
