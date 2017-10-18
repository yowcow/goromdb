package store

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWatcher_with_valid_md5sum(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "store-test")
	if err != nil {
		t.Fatal(err)
	}

	md5file, err := os.OpenFile(file.Name()+".md5", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		os.Remove(file.Name())
		os.Remove(md5file.Name())
	}()

	d := 100 * time.Millisecond
	update := make(chan bool)
	quit := make(chan bool)
	wg := &sync.WaitGroup{}

	buf := new(bytes.Buffer)
	logger := log.New(buf, "", log.LstdFlags)

	wg.Add(1)
	go NewWatcher(d, file.Name(), logger, update, quit, wg, CheckMD5Sum)

	done := make(chan bool)
	go func(done chan<- bool) {
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				file.WriteString("hogefuga\n")
				file.Close()

				md5file.WriteString("56bde24b2b0fd23d0b032c8aa128a86c  store-test")
				md5file.Close()

				timer.Reset(d)
			case <-update:
				if !timer.Stop() {
					<-timer.C
				}
				quit <- true
				wg.Wait()
				done <- true
			}
		}
	}(done)

	assert.True(t, <-done)
}

func TestKeyNotFoundError(t *testing.T) {
	err := KeyNotFoundError([]byte("hogefuga"))

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "key 'hogefuga' not found")
}
