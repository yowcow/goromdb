package simplestore

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goromdb/loader"
	"github.com/yowcow/goromdb/storage/jsonstorage"
	"github.com/yowcow/goromdb/testutil"
)

var sampleDataFile = "../../storage/jsonstorage/valid.json"

func TestNew(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	ldr, _ := loader.New(dir, "test.data")

	stg := jsonstorage.New(false)

	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	_ = New(filein, ldr, stg, logger)
}

func TestLoadAndGet(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	ldr, _ := loader.New(dir, "test.data")

	stg := jsonstorage.New(false)

	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	str := New(filein, ldr, stg, logger)
	err := str.Load(sampleDataFile)

	assert.Nil(t, err)

	key, val, err := str.Get([]byte("hoge"))

	assert.Nil(t, err)
	assert.Equal(t, []byte("hoge"), key)
	assert.Equal(t, []byte("hogehoge"), val)
}

func TestStart(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	filein := make(chan string)
	ldr, _ := loader.New(dir, "test.data")

	stg := jsonstorage.New(false)

	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	str := New(filein, ldr, stg, logger)
	done := str.Start()

	file := filepath.Join(dir, "dropin.db")
	for i := 0; i < 10; i++ {
		testutil.CopyFile(file, sampleDataFile)
		filein <- file
	}

	key, val, err := str.Get([]byte("hoge"))

	assert.Nil(t, err)
	assert.Equal(t, []byte("hoge"), key)
	assert.Equal(t, []byte("hogehoge"), val)

	close(filein)
	<-done
}
