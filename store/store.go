package store

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type Store interface {
	Get([]byte) ([]byte, error)
	Shutdown() error
}

type ChecksumFunc func(string, string) error

func NewWatcher(d time.Duration, file string, logger *log.Logger, update chan<- bool, quit <-chan bool, wg *sync.WaitGroup, checksum ChecksumFunc) {
	defer wg.Done()

	var lastModified time.Time
	md5file := file + ".md5"

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
					if err = checksum(file, md5file); err == nil {
						lastModified = fi.ModTime()
						update <- true
					} else {
						logger.Print("-> watcher file MD5 check failed: ", err)
					}
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

func CheckMD5Sum(file, md5file string) error {
	md5fh, err := os.Open(md5file)
	if err != nil {
		return err
	}

	fh, err := os.Open(file)
	if err != nil {
		return err
	}

	defer func() {
		fh.Close()
		md5fh.Close()
	}()

	expected := make([]byte, 32)
	_, err = md5fh.Read(expected)
	if err != nil {
		return err
	}

	h := md5.New()
	if _, err := io.Copy(h, fh); err != nil {
		return err
	}

	md5sum := hex.EncodeToString(h.Sum(nil))
	if md5sum != string(expected) {
		return fmt.Errorf("expecting MD5 sum '%s' but got '%s'", expected, md5sum)
	}

	return nil
}

func KeyNotFoundError(key []byte) error {
	return fmt.Errorf("key '%s' not found", string(key))
}
