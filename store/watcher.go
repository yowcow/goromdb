package store

import (
	"log"
	"os"
	"sync"
	"time"
)

type Watcher struct {
	file         string
	duration     time.Duration
	lastModified time.Time
	checksum     ChecksumChecker
	logger       *log.Logger
}

func NewWatcher(
	file string,
	duration time.Duration,
	checksum ChecksumChecker,
	logger *log.Logger,
) *Watcher {
	var lastModified time.Time
	return &Watcher{
		file,
		duration,
		lastModified,
		checksum,
		logger,
	}
}

func (w Watcher) Start(update chan<- bool, quit <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	timer := time.NewTimer(w.duration)
	w.logger.Print("-> watcher started!")

	for {
		select {
		case <-timer.C:
			if w.IsLoadable() {
				update <- true
			}
			timer.Reset(w.duration)
		case <-quit:
			if !timer.Stop() {
				<-timer.C
			}
			w.logger.Print("-> watcher finished!")
			return
		}
	}
}

func (w *Watcher) IsLoadable() bool {
	if fi, err := os.Stat(w.file); err == nil {
		if fi.ModTime() != w.lastModified {
			if w.checksum == nil {
				w.lastModified = fi.ModTime()
				return true
			} else {
				w.logger.Print("-> watcher file checksum verification in progress")
				if err = w.checksum(w.file); err == nil {
					w.logger.Print("-> watcher file checksum verification succeeded")
					w.lastModified = fi.ModTime()
					return true
				} else {
					w.logger.Print("-> watcher file checksum verification failed: ", err)
					return false
				}
			}
		}
	}
	return false
}
