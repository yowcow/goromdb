package radixhandler

import (
	"log"
	"sync"

	"github.com/armon/go-radix"
	"github.com/yowcow/goromdb/handler"
	"github.com/yowcow/goromdb/loader"
	"github.com/yowcow/goromdb/storage"
)

type Handler struct {
	tree    *radix.Tree
	storage storage.Storage
	mux     *sync.RWMutex
	logger  *log.Logger
}

func New(stg storage.Storage, logger *log.Logger) handler.Handler {
	return &Handler{
		radix.New(),
		stg,
		new(sync.RWMutex),
		logger,
	}
}

func (h *Handler) Start(filein <-chan string, l *loader.Loader) <-chan bool {
	done := make(chan bool)
	go h.start(filein, l, done)
	return done
}

func (h *Handler) start(filein <-chan string, l *loader.Loader, done chan<- bool) {
	defer func() {
		h.logger.Print("radixhandler finished")
		close(done)
	}()
	h.logger.Println("radixhandler started")
	if newfile, ok := l.FindAny(); ok {
		if err := h.Load(newfile); err != nil {
			h.logger.Printf("radixhandler failed loading data from '%s': %s", newfile, err.Error())
		}
	}
	for file := range filein {
		h.logger.Printf("radixhandler got a new file to load at '%s'", file)
		newfile, err := l.DropIn(file)
		if err != nil {
			h.logger.Printf("radixhandler failed dropping file from '%s' into '%s': %s", file, newfile, err.Error())
			continue
		}

		h.logger.Printf("radixhandler loading data from '%s'", newfile)
		err = h.Load(newfile)
		if err != nil {
			h.logger.Printf("radixhandler failed loading data from '%s': %s", newfile, err.Error())
			continue
		}

		h.logger.Printf("radixhandler successfully loaded data from '%s'", newfile)
		if ok := l.CleanUp(); ok {
			h.logger.Print("radixhandler successfully removed previously loaded file")
		}
	}
}

func (h *Handler) Load(file string) error {
	h.mux.Lock()
	err := h.storage.Load(file)
	h.mux.Unlock()
	if err != nil {
		return err
	}
	h.tree = h.buildTree()
	return nil
}

func (h Handler) buildTree() *radix.Tree {
	tree := radix.New()

	c, err := h.storage.Cursor()
	if err != nil {
		h.logger.Printf("radixhandler failed creating a tree: %s", err.Error())
		return tree
	}

	count := 0
	for {
		k, _, err := c.Next()
		if err != nil {
			h.logger.Printf("radixhandler successfully created a tree with %d keys", count)
			return tree
		}
		tree.Insert(string(k), true)
		count++
	}
}

func (h Handler) Get(key []byte) ([]byte, []byte, error) {
	h.mux.RLock()
	defer h.mux.RUnlock()

	prefix, _, ok := h.tree.LongestPrefix(string(key))
	if !ok {
		return nil, nil, storage.KeyNotFoundError(key)
	}

	p := []byte(prefix)
	val, err := h.storage.Get(p)
	if err != nil {
		return nil, nil, err
	}
	return p, val, nil
}
