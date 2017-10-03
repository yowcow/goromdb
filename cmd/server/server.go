package main

import (
	"flag"

	"github.com/yowcow/go-romdb/protocol/memcachedprotocol"
	"github.com/yowcow/go-romdb/server"
	"github.com/yowcow/go-romdb/store/jsonstore"
)

func main() {
	var addr string
	var file string
	flag.StringVar(&addr, "addr", ":11211", "Address to bind to")
	flag.StringVar(&file, "json", "./data/sample-data.json", "JSON data file")
	flag.Parse()

	proto, err := memcachedprotocol.New()

	if err != nil {
		panic(err)
	}

	store, err := jsonstore.New(file)

	if err != nil {
		panic(err)
	}

	s := server.New("tcp", addr, proto, store)
	s.Start()
}
