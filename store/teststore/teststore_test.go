package teststore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	s := New()

	assert.NotNil(t, s)
}

func TestGet_on_existing_key(t *testing.T) {
	store := New()
	v, err := store.Get([]byte("foo"))

	assert.Nil(t, err)
	assert.Equal(t, "my test foo", string(v))
}

func TestGet_on_non_existing_key(t *testing.T) {
	store := New()
	v, err := store.Get([]byte("hogefuga"))

	assert.NotNil(t, err)
	assert.Nil(t, v)
}
