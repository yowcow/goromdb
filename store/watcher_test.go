package store

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var sampleDBFile = "../data/store/sample-bdb.db"

func TestNewWatcher_with_invalid_file(t *testing.T) {
	dur := 100 * time.Millisecond
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)

	w := NewWatcher("hoge", dur, nil, logger)

	assert.NotNil(t, w)
}

func TestNewWatcher_with_valid_file(t *testing.T) {
	dur := 100 * time.Millisecond
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)
	checksum := CheckMD5Sum

	w := NewWatcher(sampleDBFile, dur, checksum, logger)

	assert.NotNil(t, w)
}

func TestWatcher_without_checksum(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "watcher-test")
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	datafile := filepath.Join(dir, "data.db")
	dur := 100 * time.Millisecond
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)

	w := NewWatcher(datafile, dur, nil, logger)
	update := make(chan bool)
	quit := make(chan bool)
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go w.Start(update, quit, wg)

	reader, err := os.Open(sampleDBFile)
	if err != nil {
		t.Fatal(err)
	}

	writer, err := os.OpenFile(datafile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		t.Fatal(err)
	}

	writer.Close()
	reader.Close()

	assert.True(t, <-update)

	quit <- true
	wg.Wait()
}

func TestWatcher_with_invalid_checksum(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "watcher-test")
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	datafile := filepath.Join(dir, "data.db")
	md5file := datafile + ".md5"
	dur := 100 * time.Millisecond
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)

	w := NewWatcher(datafile, dur, CheckMD5Sum, logger)
	update := make(chan bool)
	quit := make(chan bool)
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go w.Start(update, quit, wg)

	reader, err := os.Open(sampleDBFile)
	if err != nil {
		t.Fatal(err)
	}

	writer, err := os.OpenFile(datafile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		t.Fatal(err)
	}

	writer.Close()
	reader.Close()

	writer, err = os.OpenFile(md5file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}
	writer.WriteString("hogefugahogefugahogefugahogefuga hogefuga")
	writer.Close()

	timer := time.NewTimer(500 * time.Millisecond)
	<-timer.C
	timer.Stop()

	quit <- true
	wg.Wait()

	re := regexp.MustCompile("watcher checksum verification failed")

	assert.True(t, re.Match(buf.Bytes()))
}

func TestWatcher_with_valid_checksum(t *testing.T) {
	dir, err := ioutil.TempDir(os.TempDir(), "watcher-test")
	if err != nil {
		t.Fatal(err)
	}
	err = os.MkdirAll(dir, os.ModeDir)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	datafile := filepath.Join(dir, "data.db")
	md5file := datafile + ".md5"
	dur := 100 * time.Millisecond
	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)

	w := NewWatcher(datafile, dur, CheckMD5Sum, logger)
	update := make(chan bool)
	quit := make(chan bool)
	wg := new(sync.WaitGroup)

	wg.Add(1)
	go w.Start(update, quit, wg)

	reader, err := os.Open(sampleDBFile)
	if err != nil {
		t.Fatal(err)
	}

	writer, err := os.OpenFile(datafile, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		t.Fatal(err)
	}

	writer.Close()
	reader.Close()

	reader, err = os.Open(sampleDBFile + ".md5")
	if err != nil {
		t.Fatal(err)
	}

	writer, err = os.OpenFile(md5file, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		t.Fatal(err)
	}

	writer.Close()
	reader.Close()

	assert.True(t, <-update)

	quit <- true
	wg.Wait()

	re := regexp.MustCompile("watcher checksum verification succeeded")

	assert.True(t, re.Match(buf.Bytes()))
}
