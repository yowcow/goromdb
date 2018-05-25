package simplehandler

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

var sampleDataFile = "../../data/store/sample-data.json"

func TestNew(t *testing.T) {
	stg := jsonstorage.New(false)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	_ = New(stg, logger)
}

func TestLoadAndGet(t *testing.T) {
	stg := jsonstorage.New(false)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	h := New(stg, logger)
	err := h.Load(sampleDataFile)

	assert.Nil(t, err)

	val, err := h.Get([]byte("hoge"))

	assert.Nil(t, err)
	assert.Equal(t, []byte("hoge!"), val)
}

func TestStart(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	stg := jsonstorage.New(false)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	h := New(stg, logger)
	filein := make(chan string)
	l, _ := loader.New(dir, "test.data")
	done := h.Start(filein, l)

	file := filepath.Join(dir, "dropin.db")
	for i := 0; i < 10; i++ {
		testutil.CopyFile(file, sampleDataFile)
		filein <- file
	}

	val, err := h.Get([]byte("hoge"))

	assert.Nil(t, err)
	assert.Equal(t, []byte("hoge!"), val)

	close(filein)
	<-done
}
