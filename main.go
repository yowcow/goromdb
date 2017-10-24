package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yowcow/go-romdb/protocol"
	"github.com/yowcow/go-romdb/protocol/memcachedprotocol"
	"github.com/yowcow/go-romdb/server"
	"github.com/yowcow/go-romdb/store"
	"github.com/yowcow/go-romdb/store/bdbstore"
	"github.com/yowcow/go-romdb/store/jsonstore"
	memcachedb_bdb "github.com/yowcow/go-romdb/store/memcachedb/bdbstore"
)

func main() {
	var addr string
	var protoBackend string
	var storeBackend string
	var file string

	flag.StringVar(&addr, "addr", ":11211", "Address to bind to")
	flag.StringVar(&protoBackend, "proto", "memcached", "Protocol: memcached")
	flag.StringVar(&storeBackend, "store", "bdb", "Store: json, bdb, memcachedb-bdb")
	flag.StringVar(&file, "file", "./data/sample-bdb.db", "Data file")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	proto, err := createProtocol(protoBackend)
	if err != nil {
		panic(err)
	}

	store, err := createStore(storeBackend, file, logger)
	if err != nil {
		panic(err)
	}

	logger.Print(
		fmt.Sprintf(
			"Booting romdb server listening to address %s that talks %s protocol, with backend store %s at file path %s",
			addr, protoBackend, storeBackend, file,
		),
	)

	s := server.New("tcp", addr, proto, store, logger)
	s.Start()
}

func createProtocol(protoBackend string) (protocol.Protocol, error) {
	switch protoBackend {
	case "memcached":
		return memcachedprotocol.New(), nil
	default:
		return nil, fmt.Errorf("don't know how to handle protoc '%s'", protoBackend)
	}
}

func createStore(storeBackend, file string, logger *log.Logger) (store.Store, error) {
	switch storeBackend {
	case "bdb":
		return bdbstore.New(file, logger), nil
	case "json":
		return jsonstore.New(file, logger), nil
	case "memcachedb-bdb":
		return memcachedb_bdb.New(file, logger), nil
	default:
		return nil, fmt.Errorf("don't know how to handle store '%s'", storeBackend)
	}
}
