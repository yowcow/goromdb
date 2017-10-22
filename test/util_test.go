package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateStoreDir(t *testing.T) {
	dir, err := CreateStoreDir()
	defer os.RemoveAll(dir)

	fi, err := os.Stat(dir)

	assert.Nil(t, err)
	assert.True(t, fi.IsDir())
}
