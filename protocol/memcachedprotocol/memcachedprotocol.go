package memcachedprotocol

import (
	"bytes"
	"fmt"
	"io"

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

func (p Protocol) Reply(w io.Writer, k, v []byte) {
	fmt.Fprintf(
		w,
		"VALUE %s 0 %d\r\n%s\r\n",
		string(k),
		len(v),
		string(v),
	)
}

func (p Protocol) Finish(w io.Writer) {
	fmt.Fprint(w, "END\r\n")
}
