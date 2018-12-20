package watcher

import (
	"bytes"
	"context"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"

	"github.com/yowcow/goromdb/testutil"
)

func TestNewSimpleWatcher(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	var logbuf bytes.Buffer
	logger := log.New(&logbuf, "", 0)
	file := filepath.Join(dir, "hoge.txt")

	_ = NewSimpleWatcher(file, 1, logger)
}

func TestStartSimpleWatcher(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	var logbuf bytes.Buffer
	logger := log.New(&logbuf, "", 0)
	file := filepath.Join(dir, "hoge.txt")

	ctx, cancel := context.WithCancel(context.Background())
	w := NewSimpleWatcher(file, 1, logger)
	out := w.Start(ctx)

	time.Sleep(100 * time.Millisecond)

	cancel() // should close chan
	<-out
}

func TestSimpleWatcherVerifyFailsWhenDirIsGiven(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	var logbuf bytes.Buffer
	logger := log.New(&logbuf, "", 0)
	file := filepath.Join(dir, "hoge.txt")

	ctx, cancel := context.WithCancel(context.Background())
	w := NewSimpleWatcher(file, 10, logger)
	out := w.Start(ctx)

	err := os.Mkdir(file, 0755)
	if err != nil {
		t.Fatal("failed creating a dir", err)
	}
	time.Sleep(100 * time.Millisecond)

	cancel()
	<-out

	re := regexp.MustCompile("expected a file but got a directory")
	if re.MatchString(logbuf.String()) != true {
		t.Error("'expected a file but got ...' but got", logbuf.String())
	}
}

func TestSimpleWatcherVerifySucceeds(t *testing.T) {
	dir := testutil.CreateTmpDir()
	defer os.RemoveAll(dir)

	var logbuf bytes.Buffer
	logger := log.New(&logbuf, "", 0)
	file := filepath.Join(dir, "hoge.txt")

	ctx, cancel := context.WithCancel(context.Background())
	w := NewSimpleWatcher(file, 10, logger)
	out := w.Start(ctx)

	testutil.CopyFile(file, "valid.txt")
	time.Sleep(100 * time.Millisecond)

	fileOut := <-out

	cancel() // should close chan
	<-out

	if fileOut != file {
		t.Error("expected", file, "but got", fileOut)
	}
}
