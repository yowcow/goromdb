package bdbstore

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleDBFile = "../../../data/store/sample-memcachedb-bdb.db"

func TestNew(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.Lshortfile)

	store := New(sampleDBFile, logger)

	assert.Nil(t, store.Shutdown())
}

func TestNew_with_non_existing_file(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.Lshortfile)

	store := New("hogefuga.txt", logger)

	assert.NotNil(t, store)
	assert.Nil(t, store.Shutdown())
}

func TestGet_on_existing_key(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.Lshortfile)

	store := New(sampleDBFile, logger)
	val, err := store.Get([]byte("fuga"))

	assert.Nil(t, err)
	assert.Equal(t, "fuga!!", string(val))
	assert.Nil(t, store.Shutdown())
}

func TestGet_on_non_existing_key(t *testing.T) {
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.Lshortfile)

	store := New(sampleDBFile, logger)
	val, err := store.Get([]byte("hogefuga"))

	assert.Nil(t, val)
	assert.NotNil(t, err)
	assert.Nil(t, store.Shutdown())
}
