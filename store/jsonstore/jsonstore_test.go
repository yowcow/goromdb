package jsonstore

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	_, err := New(filein, false, logger)

	assert.Nil(t, err)
}

func TestLoad(t *testing.T) {
	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	s, _ := New(filein, false, logger)

	type Case struct {
		input       string
		expectError bool
		subtest     string
	}
	cases := []Case{
		{"invalid.json", true, "invalid json"},
		{"valid.json", false, "valid json"},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			err := s.Load(c.input)

			assert.Equal(t, c.expectError, err != nil)
		})
	}
}

func TestGet(t *testing.T) {
	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	s, _ := New(filein, false, logger)
	_ = s.Load("valid.json")

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

func TestStart(t *testing.T) {
	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	s, _ := New(filein, false, logger)
	done := s.Start()

	for i := 0; i < 10; i++ {
		filein <- "valid.json"
	}

	close(filein)
	<-done
}
