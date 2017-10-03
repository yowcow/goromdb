package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyNotFoundError(t *testing.T) {
	err := KeyNotFoundError("hogefuga")

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "key 'hogefuga' not found")
}
