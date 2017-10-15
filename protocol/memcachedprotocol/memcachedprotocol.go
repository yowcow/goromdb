package memcachedprotocol

import (
	"bufio"
	"bytes"
	"strconv"

	"github.com/yowcow/go-romdb/protocol"
)

var prefixes = [][]byte{[]byte("gets "), []byte("get ")}
var space = []byte(" ")

type MemcachedProtocol struct {
}

func New() (protocol.Protocol, error) {
	return &MemcachedProtocol{}, nil
}

func (p MemcachedProtocol) Parse(line []byte) ([][]byte, error) {
	for _, prefix := range prefixes {
		if bytes.HasPrefix(line, prefix) {
			words := bytes.Split(line, space)
			return words[1:], nil
		}
	}
	return [][]byte{}, protocol.InvalidCommandError(line)
}

func (p MemcachedProtocol) Reply(w *bufio.Writer, k, v []byte) {
	w.WriteString("VALUE ")
	w.Write(k)
	w.WriteString(" 0 ")
	w.WriteString(strconv.Itoa(len(v)))
	w.WriteString("\r\n")
	w.Write(v)
	w.WriteString("\r\n")
}

func (p MemcachedProtocol) Finish(w *bufio.Writer) {
	w.WriteString("END\r\n")
}
