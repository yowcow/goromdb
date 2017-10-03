package main

import (
	"flag"

	"github.com/yowcow/go-romdb/proto/memcachedproto"
	"github.com/yowcow/go-romdb/server"
	"github.com/yowcow/go-romdb/store/storetest"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":11211", "Address to bind to")
	flag.Parse()

	proto := memcachedproto.New()
	store := storetest.New()

	s := server.New("tcp", addr, proto, store)
	s.Start()
}
