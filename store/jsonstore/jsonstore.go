package jsonstore

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/yowcow/go-romdb/store"
)

type Data map[string]string

type JSONStore struct {
	in     chan *string
	out    chan *string
	quit   chan bool
	wg     *sync.WaitGroup
	logger *log.Logger
}

func New(file string) (store.Store, error) {
	in := make(chan *string)
	out := make(chan *string)
	boot := make(chan bool)
	quit := make(chan bool)
	wg := new(sync.WaitGroup)
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	s := &JSONStore{in, out, quit, wg, logger}

	wg.Add(1)
	go s.startDataNode(file, boot)
	<-boot // await store to boot
	s.logger.Print("-> store started!")

	return s, nil
}

func (s JSONStore) Get(key string) (string, error) {
	s.in <- &key
	v := <-s.out
	if v == nil {
		return "", store.KeyNotFoundError(key)
	}
	return *v, nil
}

func (s JSONStore) Shutdown() error {
	s.quit <- true
	s.wg.Wait()

	close(s.quit)
	close(s.in)
	close(s.out)
	s.logger.Print("-> store finished!")

	return nil
}

func (s JSONStore) startDataNode(file string, boot chan bool) {
	defer s.wg.Done()

	var data Data
	data, err := LoadJSON(file)

	if err != nil {
		s.logger.Println("-> store initial data load failed: ", err)
	}

	watcherupdate := make(chan Data)
	watcherboot := make(chan bool)
	watcherquit := make(chan bool)
	watcherwg := new(sync.WaitGroup)

	watcherwg.Add(1)
	go s.startWatcher(file, watcherupdate, watcherboot, watcherquit, watcherwg)

	<-watcherboot // await watcher to boot
	close(boot)
	s.logger.Print("-> datanode started!")

	for {
		select {
		case d := <-watcherupdate:
			data = d
			s.logger.Print("-> datanode loaded with new data!")
		case k := <-s.in:
			if v, ok := data[*k]; ok {
				s.out <- &v
			} else {
				s.out <- nil
			}
		case <-s.quit:
			watcherquit <- true
			watcherwg.Wait()

			close(watcherupdate)
			close(watcherquit)

			s.logger.Print("-> datanode finished!")
			return
		}
	}
}

func (s JSONStore) startWatcher(file string, out chan Data, boot chan bool, quit chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	var lastModified time.Time

	if fi, err := os.Stat(file); err == nil {
		lastModified = fi.ModTime()
	}

	d := 5 * time.Second
	t := time.NewTimer(d)

	close(boot)
	s.logger.Print("-> watcher started!")

	for {
		select {
		case <-t.C:
			if fi, err := os.Stat(file); err == nil {
				if fi.ModTime() != lastModified {
					lastModified = fi.ModTime()
					if data, err := LoadJSON(file); err == nil {
						out <- data
					} else {
						s.logger.Print("-> watcher failed reading data from file: ", err)
					}
				}
			} else {
				s.logger.Print("-> watcher file check failed: ", err)
			}
			t.Reset(d)
		case <-quit:
			if !t.Stop() {
				<-t.C
			}
			s.logger.Print("-> watcher finished!")
			return
		}
	}
}

func LoadJSON(file string) (Data, error) {
	var data Data

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
