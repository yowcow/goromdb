package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yowcow/goromdb/protocol"
	"github.com/yowcow/goromdb/protocol/memcachedprotocol"
	"github.com/yowcow/goromdb/server"
	"github.com/yowcow/goromdb/store"
	"github.com/yowcow/goromdb/store/bdbstore"
	"github.com/yowcow/goromdb/store/jsonstore"
	memcachedb_bdb "github.com/yowcow/goromdb/store/memcachedb/bdbstore"
)

var Version string

func main() {
	var addr string
	var protoBackend string
	var storeBackend string
	var file string
	var version bool

	flag.StringVar(&addr, "addr", ":11211", "address to bind to")
	flag.StringVar(&protoBackend, "proto", "memcached", "Protocol: memcached")
	flag.StringVar(&storeBackend, "store", "memcachedb-bdb", "Store: json, bdb, memcachedb-bdb")
	flag.StringVar(&file, "file", "/tmp/goromdb", "data file")
	flag.BoolVar(&version, "version", false, "print version")
	flag.Parse()

	if version {
		fmt.Println("goromdb", Version)
		return
	}

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
