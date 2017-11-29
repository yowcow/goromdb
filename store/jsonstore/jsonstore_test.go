package jsonstore

import (
	"bytes"
	"context"
	_ "fmt"
	"log"
	_ "os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	_, err := New(filein, logger)

	assert.Nil(t, err)
}

func TestStart(t *testing.T) {
	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	s, _ := New(filein, logger)

	ctx, cancel := context.WithCancel(context.Background())
	done := s.Start(ctx)
	cancel()
	<-done
}

func TestLoadInvalidData(t *testing.T) {
	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	s, _ := New(filein, logger)

	err := s.Load("invalid.json")

	assert.NotNil(t, err)
}

func TestLoadValidData(t *testing.T) {
	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	s, _ := New(filein, logger)
	err := s.Load("valid.json")

	assert.Nil(t, err)

	type Case struct {
		input, expected []byte
		subtest         string
	}
	cases := []Case{
		{[]byte("foo"), nil, "non-existing key"},
		{[]byte("hoge"), []byte("hogehoge"), "existing key: hoge"},
		{[]byte("fuga"), []byte("fugafuga"), "existing key: fuga"},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			actual, err := s.Get(c.input)

			if c.expected == nil {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, string(actual), string(c.expected))
			}
		})
	}
}
