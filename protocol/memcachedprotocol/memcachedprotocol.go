package memcachedprotocol

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"

	"github.com/yowcow/go-romdb/protocol"
)

var prefix = []string{"gets ", "get "}
var space = []byte(" ")

type MemcachedProtocol struct {
}

func New() (protocol.Protocol, error) {
	return &MemcachedProtocol{}, nil
}

func (p MemcachedProtocol) Parse(line []byte) ([][]byte, error) {
	for _, p := range prefix {
		if bytes.HasPrefix(line, []byte(p)) {
			words := bytes.Split(line, space)
			return words[1:], nil
		}
	}
	return [][]byte{}, fmt.Errorf("invalid command: %s", string(line))
}

func (p MemcachedProtocol) Reply(w *bufio.Writer, key string, data string) {
	w.WriteString("VALUE ")
	w.WriteString(key)
	w.WriteString(" 0 ")
	w.WriteString(strconv.Itoa(len(data)))
	w.WriteString("\r\n")
	w.WriteString(data)
	w.WriteString("\r\n")
}

func (p MemcachedProtocol) Finish(w *bufio.Writer) {
	w.WriteString("END\r\n")
}
