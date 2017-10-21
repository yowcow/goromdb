package store

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildDirs_on_existing_dir(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "store-test")
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dirs, err := BuildDirs(dir, 2)

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

func TestBuildDirs_on_non_existing_dir(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "store-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	dirs, err := BuildDirs(dir, 2)

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

func TestNewLoader(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "store-test")
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file := dir + "/hoge.txt"
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)
	loader := NewLoader(file, logger)

	assert.NotNil(t, loader)
}

func TestLoader_BuildStoreDirs(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "store-test")
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	file := dir + "/hoge.txt"
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)
	loader := NewLoader(file, logger)

	err = loader.BuildStoreDirs()

	assert.Nil(t, err)
}
