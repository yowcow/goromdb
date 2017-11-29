package watcher

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func createTmpDir() string {
	dir, err := ioutil.TempDir(os.TempDir(), "watcher")
	if err != nil {
		panic(err)
	}
	return dir
}

func copyFile(dst, src string) error {
	fo, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer fo.Close()

	fi, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fi.Close()

	if _, err = io.Copy(fo, fi); err != nil {
		return err
	}

	return nil
}

func TestVerifyFile(t *testing.T) {
	type Case struct {
		file, md5file string
		expectedOK    bool
		expectError   bool
		subtest       string
	}
	cases := []Case{
		{"non-existing.txt", "non-existing.txt.md5", false, false, "non-existing file"},
		{"valid.txt", "non-existing.txt.md5", false, false, "non-existing md5 file"},
		{"valid.txt", "invalid-len.txt.md5", false, true, "invalid md5 length"},
		{"valid.txt", "invalid-sum.txt.md5", false, true, "invalid md5 sum"},
		{"valid.txt", "valid.txt.md5", true, false, "valid md5"},
	}

	for _, c := range cases {
		t.Run(c.subtest, func(t *testing.T) {
			ok, err := verifyFile(c.file, c.md5file)

			assert.Equal(t, ok, c.expectedOK)
			assert.Equal(t, err != nil, c.expectError)
		})
	}
}

func TestNew(t *testing.T) {
	dir := createTmpDir()
	defer os.RemoveAll(dir)

	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	file := filepath.Join(dir, "hoge.txt")
	wcr := New(file, 1000, logger)

	assert.NotNil(t, wcr)
}

func TestStart(t *testing.T) {
	dir := createTmpDir()
	defer os.RemoveAll(dir)

	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	file := filepath.Join(dir, "hoge.txt")
	wcr := New(file, 1000, logger)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	out := wcr.Start(ctx)
	<-out
	cancel()
}

func TestWatchOutput(t *testing.T) {
	dir := createTmpDir()
	defer os.RemoveAll(dir)

	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	file := filepath.Join(dir, "valid.txt")
	wcr := New(file, 100, logger)

	ctx, cancel := context.WithCancel(context.Background())
	out := wcr.Start(ctx)

	err := copyFile(filepath.Join(dir, "valid.txt"), "valid.txt")
	if err != nil {
		panic(err)
	}
	err = copyFile(filepath.Join(dir, "valid.txt.md5"), "valid.txt.md5")
	if err != nil {
		panic(err)
	}

	loadedFile := <-out
	cancel()
	<-out

	assert.Equal(t, file, loadedFile)

	_, err = os.Stat(filepath.Join(dir, "valid.txt.md5"))

	assert.False(t, os.IsExist(err))
}
