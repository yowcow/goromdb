package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildStoreDirs_on_existing_dir(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "store-test")
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dir)
	defer os.RemoveAll(dir)

	dirs, err := BuildStoreDirs(dir)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(dirs))
	assert.Equal(t, filepath.Join(dir, "db00"), dirs[0])
	assert.Equal(t, filepath.Join(dir, "db01"), dirs[1])

	for _, dir := range dirs {
		fi, err := os.Stat(dir)

		assert.Nil(t, err)
		assert.True(t, fi.IsDir())
	}
}

func TestBuildStoreDirs_on_non_existing_dir(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "store-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dirs, err := BuildStoreDirs(dir)

	assert.Nil(t, err)
	assert.Equal(t, 2, len(dirs))
	assert.Equal(t, filepath.Join(dir, "db00"), dirs[0])
	assert.Equal(t, filepath.Join(dir, "db01"), dirs[1])

	for _, dir := range dirs {
		fi, err := os.Stat(dir)

		assert.Nil(t, err)
		assert.True(t, fi.IsDir())
	}
}

func TestKeyNotFoundError(t *testing.T) {
	err := KeyNotFoundError([]byte("hogefuga"))

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "key 'hogefuga' not found")
}
