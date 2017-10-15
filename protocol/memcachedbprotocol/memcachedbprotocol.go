package memcachedbprotocol

import (
	"bufio"
	"bytes"

	"github.com/yowcow/go-romdb/protocol"
	"github.com/yowcow/go-romdb/protocol/memcachedprotocol"
)

type Protocol struct {
}

func New() (protocol.Protocol, error) {
	return &Protocol{}, nil
}

func (p Protocol) Parse(line []byte) ([][]byte, error) {
	for _, prefix := range memcachedprotocol.Prefixes {
		if bytes.HasPrefix(line, prefix) {
			words := bytes.Split(line, memcachedprotocol.Space)
			return words[1:], nil
		}
	}
	return [][]byte{}, protocol.InvalidCommandError(line)
}

func (p Protocol) Reply(w *bufio.Writer, key string, data string) {
	w.WriteString(data)
}

func (p Protocol) Finish(w *bufio.Writer) {
	w.WriteString("END\r\n")
}
