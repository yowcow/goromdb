package radixgateway

import (
	"log"
	"sync"

	"github.com/armon/go-radix"
	"github.com/yowcow/goromdb/gateway"
	"github.com/yowcow/goromdb/loader"
	"github.com/yowcow/goromdb/storage"
)

type Gateway struct {
	tree    *radix.Tree
	filein  <-chan string
	loader  *loader.Loader
	storage storage.IndexableStorage
	mux     *sync.RWMutex
	logger  *log.Logger
}

func New(filein <-chan string, ldr *loader.Loader, stg storage.IndexableStorage, logger *log.Logger) gateway.Gateway {
	return &Gateway{
		radix.New(),
		filein,
		ldr,
		stg,
		new(sync.RWMutex),
		logger,
	}
}

func (g *Gateway) Start() <-chan bool {
	done := make(chan bool)
	go g.start(done)
	return done
}

func (g *Gateway) start(done chan<- bool) {
	defer func() {
		g.logger.Print("radixgateway finished")
		close(done)
	}()
	g.logger.Println("radixgateway started")
	if newfile, ok := g.loader.FindAny(); ok {
		if err := g.Load(newfile); err != nil {
			g.logger.Printf("radixgateway failed loading data from '%s': %s", newfile, err.Error())
		}
	}
	for file := range g.filein {
		g.logger.Printf("radixgateway got a new file to load at '%s'", file)
		newfile, err := g.loader.DropIn(file)
		if err != nil {
			g.logger.Printf("radixgateway failed dropping file from '%s' into '%s': %s", file, newfile, err.Error())
			continue
		}

		g.logger.Printf("radixgateway loading data from '%s'", newfile)
		err = g.Load(newfile)
		if err != nil {
			g.logger.Printf("radixgateway failed loading data from '%s': %s", newfile, err.Error())
			continue
		}

		g.logger.Printf("radixgateway successfully loaded data from '%s'", newfile)
		if ok := g.loader.CleanUp(); ok {
			g.logger.Print("radixgateway successfully removed previously loaded file")
		}
	}
}

func (g *Gateway) Load(file string) error {
	g.mux.Lock()
	err := g.storage.Load(file)
	g.mux.Unlock()
	if err != nil {
		return err
	}
	g.tree = g.buildTree()
	return nil
}

func (g Gateway) buildTree() *radix.Tree {
	tree := radix.New()
	for _, key := range g.storage.AllKeys() {
		tree.Insert(string(key), true)
	}
	return tree
}

func (g Gateway) Get(key []byte) ([]byte, []byte, error) {
	g.mux.RLock()
	defer g.mux.RUnlock()

	prefix, _, ok := g.tree.LongestPrefix(string(key))
	if !ok {
		return nil, nil, storage.KeyNotFoundError(key)
	}

	p := []byte(prefix)
	val, err := g.storage.Get(p)
	if err != nil {
		return nil, nil, err
	}
	return p, val, nil
}
