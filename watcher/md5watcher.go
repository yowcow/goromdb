package watcher

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var (
	_ Watcher = (*MD5Watcher)(nil)
)

// MD5Watcher represents a watcher with md5sum checking
type MD5Watcher struct {
	File     string
	md5file  string
	interval int
	logger   *log.Logger
}

// NewMD5Watcher returns a watcher
func NewMD5Watcher(file string, interval int, logger *log.Logger) *MD5Watcher {
	return &MD5Watcher{file, file + ".md5", interval, logger}
}

// Start starts a watcher goroutine, and returns a channel that outputs a filepath
func (w *MD5Watcher) Start(ctx context.Context) <-chan string {
	out := make(chan string)
	go w.watch(ctx, out)
	return out
}

func (w *MD5Watcher) watch(ctx context.Context, out chan<- string) {
	d := time.Duration(w.interval) * time.Millisecond
	tc := time.NewTicker(d)
	defer func() {
		w.logger.Printf("md5watcher finished watching for file: %s", w.File)
		tc.Stop()
		close(out)
	}()
	w.logger.Printf("md5watcher started watching for file: %s", w.File)
	for {
		select {
		case <-tc.C:
			if ok, err := verifyFileMD5(w.File, w.md5file); ok {
				os.Remove(w.md5file)
				out <- w.File
			} else if err != nil {
				w.logger.Println("md5watcher file verification failed:", err.Error())
			}
		case <-ctx.Done():
			return
		}
	}
}

func verifyFileMD5(file, md5file string) (bool, error) {
	fi, err := os.Open(file)
	if err != nil {
		return false, nil
	}
	defer fi.Close()

	md5fi, err := os.Open(md5file)
	if err != nil {
		return false, fmt.Errorf("file %s is found but %s is not found", file, md5file)
	}
	defer md5fi.Close()

	expectedMD5 := make([]byte, 32)
	l, err := md5fi.Read(expectedMD5)
	if err != nil {
		return false, err
	}
	if l != 32 {
		return false, fmt.Errorf("invalid md5 hex length: %d", l)
	}

	h := md5.New()
	if _, err := io.Copy(h, fi); err != nil {
		return false, err
	}

	actualMD5 := hex.EncodeToString(h.Sum(nil))
	if actualMD5 != string(expectedMD5) {
		return false, fmt.Errorf("invalid md5 sum: expected '%s' but got '%s'", string(expectedMD5), actualMD5)
	}

	return true, nil
}
