package jsonstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	store, err := New("./jsonstore-data.json")
	store.Shutdown()

	assert.Nil(t, err)
}

func TestGet_on_existing_key(t *testing.T) {
	store, _ := New("./jsonstore-data.json")
	value, err := store.Get("hoge")
	store.Shutdown()

	assert.Nil(t, err)
	assert.Equal(t, "hoge!!!", value)
}

func TestGet_on_non_existing_key(t *testing.T) {
	store, _ := New("./jsonstore-data.json")
	value, err := store.Get("foobar")
	store.Shutdown()

	assert.Equal(t, "", value)
	assert.NotNil(t, err)
}
