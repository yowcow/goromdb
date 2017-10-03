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
	quit := make(chan bool)
	wg := new(sync.WaitGroup)
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	s := &JSONStore{in, out, quit, wg, logger}

	wg.Add(1)
	go s.startDataNode(file)

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
	s.logger.Print("-> store shutdown complete")
	return nil
}

func (s JSONStore) startDataNode(file string) {
	defer s.wg.Done()

	var data Data
	update := make(chan Data)

	go s.loadJSON(file, update)
	data = <-update

	d := 5 * time.Second
	t := time.NewTimer(d)
	fi, _ := os.Stat(file)
	mt := fi.ModTime()
	reloading := false

	for {
		select {
		case d := <-update:
			data = d
			reloading = false
			s.logger.Print("-> store loaded with new data!")
		case <-t.C:
			fi, _ := os.Stat(file)
			if fi.ModTime() != mt && !reloading {
				mt = fi.ModTime()
				reloading = true
				s.logger.Print("-> store going to be reloaded!")
				go s.loadJSON(file, update)
			}
			t.Reset(d)
		case key := <-s.in:
			if v, ok := data[*key]; ok {
				s.out <- &v
			} else {
				s.out <- nil
			}
		case <-s.quit:
			s.logger.Print("-> store shutting down")
			if !t.Stop() {
				<-t.C
			}
			return
		}
	}
}

func (s JSONStore) loadJSON(file string, out chan Data) {
	var data Data
	defer func() {
		out <- data
	}()

	b, err := ioutil.ReadFile(file)
	if err != nil {
		s.logger.Print(err)
		return
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		s.logger.Print(err)
		return
	}
}
