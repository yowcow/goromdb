package memcachedbprotocol

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/yowcow/go-romdb/protocol"
	"github.com/yowcow/go-romdb/protocol/memcachedprotocol"
)

type Protocol struct {
	logger *log.Logger
}

func New() (protocol.Protocol, error) {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	return &Protocol{logger}, nil
}

func (p Protocol) Parse(line []byte) ([][]byte, error) {
	return memcachedprotocol.Parse(line)
}

func (p Protocol) Reply(w *bufio.Writer, k, v []byte) {
	r := bytes.NewReader(v)
	key, val, len, err := Deserialize(r)
	if err != nil {
		p.logger.Print("-> deserialization failed: ", err)
		return
	}

	w.WriteString("VALUE ")
	w.Write(key)
	w.WriteString(" 0 ")
	w.WriteString(strconv.Itoa(len))
	w.WriteString("\r\n")
	w.Write(val)
	w.WriteString("\r\n")
}

func (p Protocol) Finish(w *bufio.Writer) {
	memcachedprotocol.Finish(w)
}

const _Zero uint8 = 0

func Serialize(w io.Writer, key, val []byte) error {
	nKey := len(key)
	nBytes := len(val) + 2

	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, " %d %d\r\n", 0, len(val))

	sSuffix := buf.Bytes()
	nSuffix := len(sSuffix)

	var data = []interface{}{
		int32(nBytes),
		uint8(nSuffix),
		uint8(nKey),
		_Zero,
		_Zero,
		key,
		_Zero,
		sSuffix,
		val,
		[]byte("\r\n"),
	}

	for _, v := range data {
		var err error
		if err = binary.Write(w, binary.LittleEndian, v); err != nil {
			return err
		}
	}

	return nil
}

func Deserialize(r io.Reader) ([]byte, []byte, int, error) {
	var err error
	var (
		nBytes  int32
		nSuffix uint8
		nKey    uint8
		pad1    uint8
		pad2    uint8
	)
	var headers = []interface{}{
		&nBytes,
		&nSuffix,
		&nKey,
		&pad1,
		&pad2,
	}
	for _, v := range headers {
		err = binary.Read(r, binary.LittleEndian, v)
		if err != nil {
			return nil, nil, 0, err
		}
	}

	var (
		key     = make([]byte, nKey)
		sSuffix = make([]byte, nSuffix)
		val     = make([]byte, nBytes-2)
		pad3    uint8
	)
	var body = []interface{}{
		&key,
		&pad3,
		&sSuffix,
		&val,
	}
	for _, v := range body {
		err = binary.Read(r, binary.LittleEndian, v)
		if err != nil {
			return nil, nil, 0, err
		}
	}

	return key, val, int(nBytes - 2), nil
}
