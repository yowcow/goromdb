package memcachedproto

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"github.com/yowcow/go-romdb/proto"
)

type MemcachedProto struct {
	re *regexp.Regexp
}

func New() proto.Protocol {
	re := regexp.MustCompile(`^gets?\s`)
	return &MemcachedProto{re}
}

func (p MemcachedProto) Parse(line []byte) ([][]byte, error) {
	if p.re.Match(line) {
		line := p.re.ReplaceAll(line, []byte(""))
		return bytes.Split(line, []byte(" ")), nil
	}

	return [][]byte{}, fmt.Errorf("invalid command: %s", string(line))
}

func (p MemcachedProto) Reply(w *bufio.Writer, key string, data string) {
	w.WriteString("VALUE ")
	w.WriteString(key)
	w.WriteString(" 0 ")
	w.WriteString(strconv.Itoa(len(data)))
	w.WriteString("\r\n")
	w.WriteString(data)
	w.WriteString("\r\n")
}

func (p MemcachedProto) Finish(w *bufio.Writer) {
	w.WriteString("END\r\n")
}
