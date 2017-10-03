package main

import (
	"flag"

	"github.com/yowcow/go-romdb/protocol/memcachedprotocol"
	"github.com/yowcow/go-romdb/server"
	"github.com/yowcow/go-romdb/store/teststore"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":11211", "Address to bind to")
	flag.Parse()

	proto, err := memcachedprotocol.New()

	if err != nil {
		panic(err)
	}

	store, err := teststore.New()

	if err != nil {
		panic(err)
	}

	s := server.New("tcp", addr, proto, store)
	s.Start()
}
