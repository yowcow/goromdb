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
	err = os.MkdirAll(dir, 0755)
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
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)
	loader := NewLoader(dir, logger)

	assert.NotNil(t, loader)
}

func TestLoader_BuildStoreDirs(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "store-test")
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)
	loader := NewLoader(dir, logger)

	err = loader.BuildStoreDirs()

	assert.Nil(t, err)
}

func TestLoader_MoveFileToNextDir(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "store-test")
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)
	loader := NewLoader(dir, logger)
	loader.BuildStoreDirs()

	file := filepath.Join(dir, "hoge.txt")
	fh, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	_, err = fh.WriteString("hogefuga")
	if err != nil {
		t.Fatal(err)
	}
	fh.Close()

	nextFile, err := loader.MoveFileToNextDir(file)

	assert.Nil(t, err)
	assert.Equal(t, filepath.Join(dir, "db01", "hoge.txt"), nextFile)

	_, err = os.Stat(nextFile)

	assert.Nil(t, err)
}

func TestLoader_CleanOldDir(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "store-test")
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)
	loader := NewLoader(dir, logger)
	loader.BuildStoreDirs()

	oldFile := filepath.Join(dir, "db01", "hoge.txt")
	fh, err := os.OpenFile(oldFile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	_, err = fh.WriteString("hogehoge")
	if err != nil {
		t.Fatal(err)
	}
	fh.Close()

	err = loader.CleanOldDir(filepath.Join(dir, "hoge.txt"))

	assert.Nil(t, err)

	_, err = os.Stat(oldFile)

	assert.NotNil(t, err)
	assert.True(t, os.IsNotExist(err))
}
