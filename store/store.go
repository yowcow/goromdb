package store

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Store interface {
	Get([]byte) ([]byte, error)
	Shutdown() error
}

func NewWatcher(d time.Duration, file string, logger *log.Logger, update chan<- bool, quit <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	var lastModified time.Time

	if fi, err := os.Stat(file); err == nil {
		lastModified = fi.ModTime()
	}

	t := time.NewTimer(d)

	logger.Print("-> watcher started!")

	for {
		select {
		case <-t.C:
			if fi, err := os.Stat(file); err == nil {
				if fi.ModTime() != lastModified {
					lastModified = fi.ModTime()
					update <- true
				}
			} else {
				logger.Print("-> watcher file check failed: ", err)
			}
			t.Reset(d)
		case <-quit:
			if !t.Stop() {
				<-t.C
			}
			logger.Print("-> watcher finished!")
			return
		}
	}
}

func KeyNotFoundError(key []byte) error {
	return fmt.Errorf("key '%s' not found", string(key))
}
