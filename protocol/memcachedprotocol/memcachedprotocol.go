package memcachedprotocol

import (
	"bytes"
	"fmt"
	"io"

	"github.com/yowcow/goromdb/protocol"
)

// Prefixes defines memcached protocol command prefixes to parse
var Prefixes = [][]byte{[]byte("gets "), []byte("get ")}

// Space defines a space in []byte
var Space = []byte(" ")

// Protocol represents a protocol
type Protocol struct {
}

// New creates a new protocol
func New() protocol.Protocol {
	return &Protocol{}
}

// Parse parses given line into keys to search
func (p Protocol) Parse(line []byte) ([][]byte, error) {
	for _, prefix := range Prefixes {
		if bytes.HasPrefix(line, prefix) {
			words := bytes.Split(line, Space)
			return words[1:], nil
		}
	}
	return [][]byte{}, protocol.InvalidCommandError(line)
}

// Reply writes reply message to writer
func (p Protocol) Reply(w io.Writer, k, v []byte) {
	fmt.Fprintf(w, "VALUE %s 0 %d\r\n%s\r\n", string(k), len(v), string(v))
}

// Finish writes an end of message to writer
func (p Protocol) Finish(w io.Writer) {
	fmt.Fprint(w, "END\r\n")
}
