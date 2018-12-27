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

var sampleNSDataFile = "../../_data/store/sample-ns-data.json"

func TestNewNS(t *testing.T) {
	stg := jsonstorage.NewNS(false)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	_ = NewNS(stg, logger)
}

func TestLoadAndGetNS(t *testing.T) {
	stg := jsonstorage.NewNS(false)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	h := NewNS(stg, logger)
	err := h.Load(sampleNSDataFile)

	assert.Nil(t, err)

	val, err := h.GetNS([]byte("ns1"), []byte("hoge"))

	assert.Nil(t, err)
	assert.Equal(t, []byte("hoge1"), val)
}

func TestStartNS(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	stg := jsonstorage.NewNS(false)
	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	h := NewNS(stg, logger)
	filein := make(chan string)
	l, _ := loader.New(dir, "test.data")
	done := h.Start(filein, l)

	file := filepath.Join(dir, "dropin.db")
	for i := 0; i < 10; i++ {
		testutil.CopyFile(file, sampleNSDataFile)
		filein <- file
	}

	val, err := h.GetNS([]byte("ns2"), []byte("hoge"))

	assert.Nil(t, err)
	assert.Equal(t, []byte("hoge2"), val)

	close(filein)
	<-done
}
