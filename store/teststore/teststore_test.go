package teststore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	_, err := New()

	assert.Nil(t, err)
}

func TestGet_on_existing_key(t *testing.T) {
	store, _ := New()
	v, err := store.Get("foo")

	assert.Nil(t, err)
	assert.Equal(t, "my test foo", v)
}

func TestGet_on_non_existing_key(t *testing.T) {
	store, _ := New()
	v, err := store.Get("hogefuga")

	assert.NotNil(t, err)
	assert.Equal(t, "", v)
}
