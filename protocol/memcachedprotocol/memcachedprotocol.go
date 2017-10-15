package memcachedprotocol

import (
	"bufio"
	"bytes"
	"strconv"

	"github.com/yowcow/go-romdb/protocol"
)

var Prefixes = [][]byte{[]byte("gets "), []byte("get ")}
var Space = []byte(" ")

type Protocol struct {
}

func New() (protocol.Protocol, error) {
	return &Protocol{}, nil
}

func (p Protocol) Parse(line []byte) ([][]byte, error) {
	for _, prefix := range Prefixes {
		if bytes.HasPrefix(line, prefix) {
			words := bytes.Split(line, Space)
			return words[1:], nil
		}
	}
	return [][]byte{}, protocol.InvalidCommandError(line)
}

func (p Protocol) Reply(w *bufio.Writer, k, v []byte) {
	w.WriteString("VALUE ")
	w.Write(k)
	w.WriteString(" 0 ")
	w.WriteString(strconv.Itoa(len(v)))
	w.WriteString("\r\n")
	w.Write(v)
	w.WriteString("\r\n")
}

func (p Protocol) Finish(w *bufio.Writer) {
	w.WriteString("END\r\n")
}
