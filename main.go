package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/yowcow/goromdb/gateway"
	"github.com/yowcow/goromdb/gateway/radixgateway"
	"github.com/yowcow/goromdb/gateway/simplegateway"
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
	var gatewayBackend string
	var storageBackend string
	var file string
	var gzipped bool
	var basedir string
	var help bool
	var version bool

	flag.StringVar(&addr, "addr", ":11211", "address to bind to")
	flag.StringVar(&protoBackend, "proto", "memcached", "protocol: memcached")
	flag.StringVar(&gatewayBackend, "gateway", "simple", "gateway: simple, radix")
	flag.StringVar(&storageBackend, "storage", "json", "storage: json, bdb, bdb-memcachedb")
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

	ldr, err := loader.New(basedir, "data.db")
	if err != nil {
		panic(err)
	}

	gw, err := createGateway(gatewayBackend, filein, ldr, stg, logger)
	if err != nil {
		panic(err)
	}
	done := gw.Start()

	logger.Printf(
		"booting goromdb (address: %s, protocol: %s, gateway: %s, storage: %s, file: %s)",
		addr, protoBackend, gatewayBackend, storageBackend, file,
	)

	svr := server.New("tcp", addr, proto, gw, logger)
	err = svr.Start()
	if err != nil {
		logger.Printf("failed booting goromdb: %s", err.Error())
		os.Exit(1)
	}
	cancel()
	<-done

}

func createGateway(
	gatewayBackend string,
	filein <-chan string,
	ldr *loader.Loader,
	stg storage.IndexableStorage,
	logger *log.Logger,
) (gateway.Gateway, error) {
	switch gatewayBackend {
	case "simple":
		return simplegateway.New(filein, ldr, stg, logger), nil
	case "radix":
		return radixgateway.New(filein, ldr, stg, logger), nil
	default:
		return nil, fmt.Errorf("don't know how to handle gateway '%s'", gatewayBackend)
	}
}

func createStorage(storageBackend string, gzipped bool) (storage.IndexableStorage, error) {
	switch storageBackend {
	case "json":
		return jsonstorage.New(gzipped), nil
	case "bdb":
		return bdbstorage.New(), nil
	case "bdb-memcachedb":
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
