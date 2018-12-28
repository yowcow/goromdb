package watcher

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	_ Watcher = (*SimpleWatcher)(nil)
)

// SimpleWatcher represents a watcher without any specific checking
type SimpleWatcher struct {
	file     string
	interval int
	logger   *log.Logger
}

// NewSimpleWatcher returns a SimpleWatcher
func NewSimpleWatcher(file string, interval int, logger *log.Logger) *SimpleWatcher {
	return &SimpleWatcher{file, interval, logger}
}

// Start starts a watcher goroutine, and returns a channel that emits a filepath
func (w *SimpleWatcher) Start(ctx context.Context) <-chan string {
	out := make(chan string)
	go w.watch(ctx, out)
	return out
}

func (w *SimpleWatcher) watch(ctx context.Context, out chan<- string) {
	d := time.Duration(w.interval) * time.Millisecond
	tc := time.NewTicker(d)
	defer func() {
		w.logger.Printf("simplewatcher finished watching for file: %s", w.file)
		tc.Stop()
		close(out)
	}()
	w.logger.Printf("simplewatcher started watching for file: %s", w.file)
	for {
		select {
		case <-tc.C:
			if ok, err := verifyFileSimple(w.file); ok {
				out <- w.file
			} else if err != nil {
				w.logger.Println("simplewatcher file verification failed:", err.Error())
			}
		case <-ctx.Done():
			return
		}
	}
}

func verifyFileSimple(file string) (bool, error) {
	fi, err := os.Stat(file)
	if err != nil {
		return false, nil
	}

	fm := fi.Mode()

	if fm.IsDir() {
		return false, fmt.Errorf("expected a file but got a directory for %s", file)
	}

	if !fm.IsRegular() {
		return false, fmt.Errorf("expected a regular file but got something-else for %s", file)
	}

	return true, nil
}
