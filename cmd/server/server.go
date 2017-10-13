package main

import (
	"flag"
	"fmt"

	"github.com/yowcow/go-romdb/protocol/memcachedprotocol"
	"github.com/yowcow/go-romdb/server"
	"github.com/yowcow/go-romdb/store"
	"github.com/yowcow/go-romdb/store/bdbstore"
	"github.com/yowcow/go-romdb/store/jsonstore"
)

func main() {
	var addr string
	var storeBackend string
	var file string

	flag.StringVar(&addr, "addr", ":11211", "Address to bind to")
	flag.StringVar(&storeBackend, "store", "bdb", "Store type: json, bdb")
	flag.StringVar(&file, "file", "./data/sample-bdb.db", "Data file")
	flag.Parse()

	proto, err := memcachedprotocol.New()

	if err != nil {
		panic(err)
	}

	store, err := createStore(storeBackend, file)

	if err != nil {
		panic(err)
	}

	s := server.New("tcp", addr, proto, store)
	s.Start()
}

func createStore(storeBackend, file string) (store.Store, error) {
	switch storeBackend {
	case "bdb":
		return bdbstore.New(file)
	case "json":
		return jsonstore.New(file)
	default:
		return nil, fmt.Errorf("don't know how to handle store '%s'", storeBackend)
	}
}
