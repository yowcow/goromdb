package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yowcow/goromdb/protocol"
	"github.com/yowcow/goromdb/protocol/memcachedprotocol"
	"github.com/yowcow/goromdb/server"
	"github.com/yowcow/goromdb/store"
	"github.com/yowcow/goromdb/store/jsonstore"
	"github.com/yowcow/goromdb/watcher"
)

var Version string

func main() {
	var addr string
	var protoBackend string
	var storeBackend string
	var file string
	var version bool

	flag.StringVar(&addr, "addr", ":11211", "address to bind to")
	flag.StringVar(&protoBackend, "proto", "memcached", "protocol: memcached")
	flag.StringVar(&storeBackend, "store", "jsonstore", "store: jsonstore")
	flag.StringVar(&file, "file", "/tmp/goromdb", "data file to load into store")
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

	ctx, cancel := context.WithCancel(context.Background())
	wcr := watcher.New(file, 5000, logger)
	filein := wcr.Start(ctx)

	st, err := createStore(storeBackend, filein, logger)
	if err != nil {
		panic(err)
	}
	done := st.Start()

	logger.Printf(
		"booting goromdb (address: %s, protocol: %s, store: %s, file: %s)",
		addr, protoBackend, storeBackend, file,
	)

	svr := server.New("tcp", addr, proto, st, logger)
	err = svr.Start()
	if err != nil {
		logger.Printf("failed booting goromdb: %s", err.Error())
		os.Exit(1)
	}
	cancel()
	<-done

}

func createProtocol(protoBackend string) (protocol.Protocol, error) {
	switch protoBackend {
	case "memcached":
		return memcachedprotocol.New(), nil
	default:
		return nil, fmt.Errorf("don't know how to handle protocol '%s'", protoBackend)
	}
}

func createStore(storeBackend string, filein <-chan string, logger *log.Logger) (store.Store, error) {
	switch storeBackend {
	case "jsonstore":
		return jsonstore.New(filein, logger)
	default:
		return nil, fmt.Errorf("don't know how to handle store '%s'", storeBackend)
	}
}
