package store

import (
	"io/ioutil"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWatcher(t *testing.T) {
	file, err := ioutil.TempFile(os.TempDir(), "store-test")

	assert.Nil(t, err)

	defer os.Remove(file.Name())

	d := 100 * time.Millisecond
	update := make(chan bool)
	quit := make(chan bool)
	wg := &sync.WaitGroup{}

	logger := log.New(os.Stdout, "", log.LstdFlags)

	wg.Add(1)
	go NewWatcher(d, file.Name(), logger, update, quit, wg)

	done := make(chan bool)
	go func(done chan<- bool) {
		timer := time.NewTimer(d)
		for {
			select {
			case <-timer.C:
				file.WriteString("hoge")
				file.Sync()
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
	err := KeyNotFoundError("hogefuga")

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "key 'hogefuga' not found")
}
