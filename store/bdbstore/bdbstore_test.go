package bdbstore

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goromdb/loader"
	"github.com/yowcow/goromdb/testutil"
)

var sampleDBFile = "../../data/store/sample-bdb.db"

func TestNew(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	ldr, _ := loader.New(dir, "data.db")
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)
	_, err := New(filein, ldr, logger)

	assert.Nil(t, err)
}

func TestLoad(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	ldr, _ := loader.New(dir, "data.db")
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", 0)
	s, _ := New(filein, ldr, logger)

	type Case struct {
		input       string
		expectError bool
		subtest     string
	}
	cases := []Case{
		{
			sampleDBFile + ".hoge",
			true,
			"non-existing file fails",
		},
		{
			"../../data/store/sample-data.json",
			true,
			"non-bdb file fails",
		},
		{
			sampleDBFile,
			false,
			"existing bdb file succeeds",
		},
		{
			sampleDBFile,
			false,
			"another bdb file succeeds",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			err := s.Load(c.input)
			assert.Equal(t, c.expectError, err != nil)
		})
	}
}

func TestGet(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	ldr, _ := loader.New(dir, "data.db")
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	s, _ := New(filein, ldr, logger)
	s.Load(sampleDBFile)

	type Case struct {
		input       string
		expectedKey []byte
		expectedVal []byte
		expectError bool
		subtest     string
	}
	cases := []Case{
		{
			"hoge",
			[]byte("hoge"),
			[]byte("hoge!"),
			false,
			"existing key returns expected val",
		},
		{
			"hoge",
			[]byte("hoge"),
			[]byte("hoge!"),
			false,
			"existing key again returns expected val",
		},
		{
			"hogehoge",
			nil,
			nil,
			true,
			"non-existing key returns error",
		},
		{
			"hogehoge",
			nil,
			nil,
			true,
			"non-existing key again returns error",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			key, val, err := s.Get([]byte(c.input))

			assert.Equal(t, c.expectError, err != nil)
			assert.Equal(t, c.expectedKey, key)
			assert.Equal(t, c.expectedVal, val)
		})
	}
}

func TestStart(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	ldr, _ := loader.New(dir, "data.db")
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	s, _ := New(filein, ldr, logger)
	done := s.Start()

	file := filepath.Join(dir, "dropin.db")
	for i := 0; i < 10; i++ {
		testutil.CopyFile(file, sampleDBFile)
		filein <- file
	}

	key, val, err := s.Get([]byte("hoge"))

	assert.Nil(t, err)
	assert.Equal(t, "hoge", string(key))
	assert.Equal(t, "hoge!", string(val))

	close(filein)
	<-done
}
