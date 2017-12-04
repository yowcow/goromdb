package simplegateway

import (
	"log"
	"sync"

	"github.com/yowcow/goromdb/gateway"
	"github.com/yowcow/goromdb/loader"
	"github.com/yowcow/goromdb/storage"
)

type Gateway struct {
	filein  <-chan string
	loader  *loader.Loader
	storage storage.Storage
	mux     *sync.RWMutex
	logger  *log.Logger
}

func New(filein <-chan string, ldr *loader.Loader, stg storage.Storage, logger *log.Logger) gateway.Gateway {
	return &Gateway{
		filein,
		ldr,
		stg,
		new(sync.RWMutex),
		logger,
	}
}

func (g Gateway) Start() <-chan bool {
	done := make(chan bool)
	go g.start(done)
	return done
}

func (g Gateway) start(done chan<- bool) {
	defer func() {
		g.logger.Println("simplegateway finished")
		close(done)
	}()
	g.logger.Println("simplegateway started")
	if newfile, ok := g.loader.FindAny(); ok {
		if err := g.Load(newfile); err != nil {
			g.logger.Printf("simplegateway failed loading data from '%s': %s", newfile, err.Error())
		}
	}
	for file := range g.filein {
		g.logger.Printf("simplegateway got a new file to load at '%s'", file)
		newfile, err := g.loader.DropIn(file)
		if err != nil {
			g.logger.Printf("simplegateway failed dropping file from '%s' into '%s': %s", file, newfile, err.Error())
			continue
		}

		g.logger.Printf("simplegateway loading data from '%s'", newfile)
		err = g.Load(newfile)
		if err != nil {
			g.logger.Printf("simplegateway failed loading data from '%s': %s", newfile, err.Error())
			continue
		}

		g.logger.Printf("simplegateway successfully loaded data from '%s'", newfile)
		if ok := g.loader.CleanUp(); ok {
			g.logger.Print("simplegateway successfully removed previously loaded file")
		}
	}
}

func (g Gateway) Load(file string) error {
	g.mux.Lock()
	defer g.mux.Unlock()
	return g.storage.Load(file)
}

func (g Gateway) Get(key []byte) ([]byte, []byte, error) {
	g.mux.RLock()
	defer g.mux.RUnlock()
	val, err := g.storage.Get(key)
	if err != nil {
		return nil, nil, err
	}
	return key, val, nil
}
