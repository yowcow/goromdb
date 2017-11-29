package mdbstore

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goromdb/store/bdbstore"
	"github.com/yowcow/goromdb/testutil"
)

var sampleDBFile = "../../data/store/sample-memcachedb-bdb.db"

func TestSerialize_and_Deserialize(t *testing.T) {
	buf := new(bytes.Buffer)
	err := Serialize(buf, []byte("hoge"), []byte("ほげほげ!!"))

	assert.Nil(t, err)

	r := bytes.NewReader(buf.Bytes())
	key, val, len, err := Deserialize(r)

	assert.Nil(t, err)
	assert.Equal(t, []byte("hoge"), key)
	assert.Equal(t, []byte("ほげほげ!!"), val)
	assert.Equal(t, 14, len)
}

func TestDeserialize_returns_header_error(t *testing.T) {
	r := bytes.NewReader([]byte(""))
	_, _, _, err := Deserialize(r)

	assert.Equal(t, "failed reading memcachedb binary headers: EOF", err.Error())
}

func TestDeserialize_returns_body_error(t *testing.T) {
	r := bytes.NewReader([]byte("hogefuga"))
	_, _, _, err := Deserialize(r)

	assert.Equal(t, "failed reading memcachedb binary body: EOF", err.Error())
}

func TestNew(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	bdb, err := bdbstore.New(filein, dir, logger)

	assert.Nil(t, err)

	_, err = New(bdb, logger)

	assert.Nil(t, err)
}

func TestLoad(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	bdb, _ := bdbstore.New(filein, dir, logger)
	mdb, _ := New(bdb, logger)

	type Case struct {
		input       string
		expectError bool
		subtest     string
	}
	cases := []Case{
		{
			dir,
			true,
			"loading dir fails",
		},
		{
			sampleDBFile + ".hoge",
			true,
			"loading non-existing file fails",
		},
		{
			"../../data/store/sample-data.json",
			true,
			"loading non-bdb file fails",
		},
		{
			sampleDBFile,
			false,
			"loading bdb file succeeds",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			err := mdb.Load(c.input)
			assert.Equal(t, c.expectError, err != nil)
		})
	}
}

func TestGet(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	bdb, _ := bdbstore.New(filein, dir, logger)
	s, _ := New(bdb, logger)
	s.Load(sampleDBFile)

	type Case struct {
		input       string
		expectedVal []byte
		expectError bool
		subtest     string
	}
	cases := []Case{
		{
			"hoge",
			[]byte("hoge!"),
			false,
			"existing key returns expected val",
		},
		{
			"hoge",
			[]byte("hoge!"),
			false,
			"existing key again returns expected val",
		},
		{
			"hogehoge",
			nil,
			true,
			"non-existing key returns error",
		},
		{
			"hogehoge",
			nil,
			true,
			"non-existing key again returns error",
		},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			val, err := s.Get([]byte(c.input))

			assert.Equal(t, c.expectError, err != nil)
			assert.Equal(t, c.expectedVal, val)
		})
	}
}

func TestStart(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)
	bdb, _ := bdbstore.New(filein, dir, logger)
	s, _ := New(bdb, logger)
	s.Load(sampleDBFile)
	done := s.Start()

	file := filepath.Join(dir, "dropin.db")
	for i := 0; i < 10; i++ {
		testutil.CopyFile(file, sampleDBFile)
		filein <- file
	}

	val, err := s.Get([]byte("hoge"))

	assert.Nil(t, err)
	assert.Equal(t, "hoge!", string(val))

	close(filein)
	<-done
}
