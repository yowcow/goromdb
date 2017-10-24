package jsonstore

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goromdb/test"
)

var sampleDBFile = "../../data/store/sample-data.json"

func TestLoadJSON_returns_error_on_non_existing_file(t *testing.T) {
	_, err := LoadJSON("./hoge/fuga")

	assert.NotNil(t, err)
}

func TestLoadJSON_returns_error_on_invalid_JSON(t *testing.T) {
	_, err := LoadJSON("./jsonstore-invalid.json")

	assert.NotNil(t, err)
}

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

	assert.NotNil(t, store)
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
	store := New(filepath.Join(dir, "hogehoge.txt"), logger)

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
	value, err := store.Get([]byte("hoge"))

	assert.Nil(t, err)
	assert.Equal(t, "hoge!", string(value))
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
	value, err := store.Get([]byte("foobar"))

	assert.Nil(t, value)
	assert.NotNil(t, err)
	assert.Nil(t, store.Shutdown())
}
