package server

import (
	"bytes"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goromdb/testutil"
)

func TestHandleConn(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	sock := filepath.Join(dir, "test.sock")
	svr := New("unix", sock, logger)

	type Case struct {
		subtest      string
		input        []byte
		expectedLine []byte
	}
	cases := []Case{
		{
			"a line that end with \\r\\n",
			[]byte("hello world\r\n"),
			[]byte("hello world"),
		},
	}

	for _, c := range cases {
		done := make(chan bool)

		ln, err := net.Listen("unix", sock)
		if err != nil {
			panic(err)
		}
		go func() {
			defer close(done)
			for {
				conn, err := ln.Accept()
				if err != nil {
					return
				}
				svr.HandleConn(conn, OnReadCallbackFunc(func(conn net.Conn, line []byte, logger *log.Logger) {
					assert.Equal(t, c.expectedLine, line)
				}))
			}
		}()

		conn, err := net.Dial("unix", sock)
		if err != nil {
			panic(err)
		}

		_, err = conn.Write(c.input)

		conn.Close()
		ln.Close()
		<-done

		assert.Nil(t, err)
	}
}
