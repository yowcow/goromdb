package store

import (
	"log"
	"os"
	"sync"
	"time"
)

// Watcher represents a watcher
type Watcher struct {
	file         string
	duration     time.Duration
	lastModified time.Time
	checksum     ChecksumChecker
	logger       *log.Logger
}

// NewWatcher creates a new watcher
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

// Start watches file update, and notifies to given channel when updated
func (w Watcher) Start(update chan<- bool, quit <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	timer := time.NewTimer(w.duration)
	w.logger.Print("watcher started")
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
			w.logger.Print("watcher finished")
			return
		}
	}
}

// IsLoadable determines if the file is now loadable
func (w *Watcher) IsLoadable() bool {
	if fi, err := os.Stat(w.file); err == nil {
		if fi.ModTime() != w.lastModified {
			if w.checksum == nil {
				w.lastModified = fi.ModTime()
				return true
			}
			w.logger.Print("watcher checksum verification in progress...")
			if err = w.checksum(w.file); err == nil {
				w.logger.Print("watcher checksum verification succeeded")
				w.lastModified = fi.ModTime()
				return true
			}
			w.logger.Print("watcher checksum verification failed: ", err)
			return false
		}
	}
	return false
}
