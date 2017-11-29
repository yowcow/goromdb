package watcher

import (
	"bytes"
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yowcow/goromdb/testutil"
)

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
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	file := filepath.Join(dir, "hoge.txt")
	wcr := New(file, 1000, logger)

	assert.NotNil(t, wcr)
}

func TestStart(t *testing.T) {
	dir := testutil.CreateTmpDir()
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
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	logbuf := new(bytes.Buffer)
	logger := log.New(logbuf, "", 0)

	file := filepath.Join(dir, "valid.txt")
	wcr := New(file, 100, logger)

	ctx, cancel := context.WithCancel(context.Background())
	out := wcr.Start(ctx)

	testutil.CopyFile(filepath.Join(dir, "valid.txt"), "valid.txt")
	testutil.CopyFile(filepath.Join(dir, "valid.txt.md5"), "valid.txt.md5")

	loadedFile := <-out
	cancel()
	<-out

	assert.Equal(t, file, loadedFile)

	_, err := os.Stat(filepath.Join(dir, "valid.txt.md5"))

	assert.False(t, os.IsExist(err))
}
