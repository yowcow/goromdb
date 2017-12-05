package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yowcow/goromdb/handler"
	"github.com/yowcow/goromdb/handler/radixhandler"
	"github.com/yowcow/goromdb/handler/simplehandler"
	"github.com/yowcow/goromdb/loader"
	"github.com/yowcow/goromdb/protocol"
	"github.com/yowcow/goromdb/protocol/memcachedprotocol"
	"github.com/yowcow/goromdb/server"
	"github.com/yowcow/goromdb/storage"
	"github.com/yowcow/goromdb/storage/bdbstorage"
	"github.com/yowcow/goromdb/storage/jsonstorage"
	"github.com/yowcow/goromdb/storage/memcdstorage"
	"github.com/yowcow/goromdb/watcher"
)

var Version string

func main() {
	var addr string
	var protoBackend string
	var handlerBackend string
	var storageBackend string
	var file string
	var gzipped bool
	var basedir string
	var help bool
	var version bool

	flag.StringVar(&addr, "addr", ":11211", "address to bind to")
	flag.StringVar(&protoBackend, "proto", "memcached", "protocol: memcached")
	flag.StringVar(&handlerBackend, "handler", "simple", "handler: simple, radix")
	flag.StringVar(&storageBackend, "storage", "json", "storage: json, bdb, memcachedb-bdb")
	flag.StringVar(&file, "file", "/tmp/goromdb", "data file to be loaded into store")
	flag.BoolVar(&gzipped, "gzipped", false, "whether or not loading file is gzipped")
	flag.StringVar(&basedir, "basedir", "", "base directory to store loaded data file")
	flag.BoolVar(&help, "help", false, "print help")
	flag.BoolVar(&version, "version", false, "print version")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if version {
		fmt.Println("goromdb", Version)
		os.Exit(0)
	}

	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	proto, err := createProtocol(protoBackend)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wcr := watcher.New(file, 5000, logger)
	filein := wcr.Start(ctx)

	stg, err := createStorage(storageBackend, gzipped)
	if err != nil {
		panic(err)
	}

	l, err := loader.New(basedir, "data.db")
	if err != nil {
		panic(err)
	}

	h, err := createHandler(handlerBackend, stg, logger)
	if err != nil {
		panic(err)
	}
	done := h.Start(filein, l)

	logger.Printf(
		"booting goromdb (address: %s, protocol: %s, handler: %s, storage: %s, file: %s)",
		addr, protoBackend, handlerBackend, storageBackend, file,
	)

	svr := server.New("tcp", addr, proto, h, logger)
	err = svr.Start()
	if err != nil {
		logger.Printf("failed booting goromdb: %s", err.Error())
		os.Exit(1)
	}
	cancel()
	<-done

}

func createHandler(
	handlerBackend string,
	stg storage.IndexableStorage,
	logger *log.Logger,
) (handler.Handler, error) {
	switch handlerBackend {
	case "simple":
		return simplehandler.New(stg, logger), nil
	case "radix":
		return radixhandler.New(stg, logger), nil
	default:
		return nil, fmt.Errorf("don't know how to handle handler '%s'", handlerBackend)
	}
}

func createStorage(storageBackend string, gzipped bool) (storage.IndexableStorage, error) {
	switch storageBackend {
	case "json":
		return jsonstorage.New(gzipped), nil
	case "bdb":
		return bdbstorage.New(), nil
	case "memcachedb-bdb":
		p := bdbstorage.New()
		return memcdstorage.New(p), nil
	default:
		return nil, fmt.Errorf("don't know how to handle storage '%s'", storageBackend)
	}
}

func createProtocol(protoBackend string) (protocol.Protocol, error) {
	switch protoBackend {
	case "memcached":
		return memcachedprotocol.New(), nil
	default:
		return nil, fmt.Errorf("don't know how to handle protocol '%s'", protoBackend)
	}
}
