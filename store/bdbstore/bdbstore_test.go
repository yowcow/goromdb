package bdbstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var sampleDBFile = "../../data/sample-bdb.db"

func TestNew(t *testing.T) {
	store, err := New(sampleDBFile)

	assert.Nil(t, err)
	assert.Nil(t, store.Shutdown())
}

func TestGet_on_existing_key(t *testing.T) {
	store, _ := New(sampleDBFile)
	val, err := store.Get("fuga")

	assert.Nil(t, err)
	assert.Equal(t, "fuga!!", val)
	assert.Nil(t, store.Shutdown())
}

func TestGet_on_non_existing_key(t *testing.T) {
	store, _ := New(sampleDBFile)
	val, err := store.Get("hogefuga")

	assert.Equal(t, "", val)
	assert.NotNil(t, err)
	assert.Nil(t, store.Shutdown())
}
