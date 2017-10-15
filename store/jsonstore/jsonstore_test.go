package jsonstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadJSON_returns_error_on_non_existing_file(t *testing.T) {
	_, err := LoadJSON("./hoge/fuga")

	assert.NotNil(t, err)
}

func TestLoadJSON_returns_error_on_invalid_JSON(t *testing.T) {
	_, err := LoadJSON("./jsonstore-invalid.json")

	assert.NotNil(t, err)
}

func TestNew(t *testing.T) {
	store, err := New("./jsonstore-data.json")
	store.Shutdown()

	assert.Nil(t, err)
}

func TestGet_on_existing_key(t *testing.T) {
	store, _ := New("./jsonstore-data.json")
	value, err := store.Get([]byte("hoge"))

	assert.Nil(t, err)
	assert.Equal(t, "hoge!!!", string(value))

	store.Shutdown()
}

func TestGet_on_non_existing_key(t *testing.T) {
	store, _ := New("./jsonstore-data.json")
	value, err := store.Get([]byte("foobar"))
	store.Shutdown()

	assert.Nil(t, value)
	assert.NotNil(t, err)
}
