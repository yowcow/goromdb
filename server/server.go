package server

import (
	"bufio"
	"log"
	"net"
	"os"

	"github.com/yowcow/go-romdb/protocol"
	"github.com/yowcow/go-romdb/store"
)

type Server struct {
	proto    string
	addr     string
	protocol protocol.Protocol
	store    store.Store
	logger   *log.Logger
}

func New(proto, addr string, protocol protocol.Protocol, store store.Store) *Server {
	logger := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	return &Server{proto, addr, protocol, store, logger}
}

func (s Server) Start() error {
	ln, err := net.Listen(s.proto, s.addr)
	defer ln.Close()

	if err != nil {
		s.logger.Print(err)
		return err
	}

	s.logger.Print("server listening on port:", s.addr)

	for {
		conn, err := ln.Accept()

		if err != nil {
			s.logger.Print("-> failed accepting a new conn:", err)
		} else {
			s.logger.Print("-> accepted a new conn")
			go s.handleConn(conn)
		}
	}
}

func (s Server) handleConn(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	w := bufio.NewWriter(conn)

	for {
		line, _, err := r.ReadLine()

		if err != nil {
			s.logger.Print("-> read error:", err)
			break
		}

		s.logger.Print("-> read a line:", string(line))

		keys, err := s.protocol.Parse(line)

		if err != nil {
			s.logger.Print("-> protocol error:", err)
		} else {
			for _, k := range keys {
				key := string(k)
				if v, err := s.store.Get(key); err == nil {
					s.protocol.Reply(w, key, v)
				} else {
					s.logger.Print("-> get error:", err)
				}
			}
		}

		s.protocol.Finish(w)
		w.Flush()
	}
}
