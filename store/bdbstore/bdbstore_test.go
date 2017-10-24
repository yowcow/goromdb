package bdbstore

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goromdb/test"
)

var sampleDBFile = "../../data/store/sample-bdb.db"

func TestNew(t *testing.T) {
	dir, err := test.CreateStoreDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file, err := test.CopyDBFile(dir, sampleDBFile)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.Lshortfile)
	store := New(file, logger)

	assert.Nil(t, store.Shutdown())
}

func TestNew_with_non_existing_file(t *testing.T) {
	dir, err := test.CreateStoreDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.Lshortfile)
	store := New(filepath.Join(dir, "hogefuga.txt"), logger)

	assert.NotNil(t, store)
	assert.Nil(t, store.Shutdown())
}

func TestGet_on_existing_key(t *testing.T) {
	dir, err := test.CreateStoreDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file, err := test.CopyDBFile(dir, sampleDBFile)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.Lshortfile)
	store := New(file, logger)
	val, err := store.Get([]byte("fuga"))

	assert.Nil(t, err)
	assert.Equal(t, "fuga!!", string(val))
	assert.Nil(t, store.Shutdown())
}

func TestGet_on_non_existing_key(t *testing.T) {
	dir, err := test.CreateStoreDir()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file, err := test.CopyDBFile(dir, sampleDBFile)
	if err != nil {
		t.Fatal(err)
	}

	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.Lshortfile)
	store := New(file, logger)
	val, err := store.Get([]byte("hogefuga"))

	assert.Nil(t, val)
	assert.NotNil(t, err)
	assert.Nil(t, store.Shutdown())
}
