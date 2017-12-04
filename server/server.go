package server

import (
	"bufio"
	"io"
	"log"
	"net"

	"github.com/yowcow/goromdb/gateway"
	"github.com/yowcow/goromdb/protocol"
)

// Server represents a server
type Server struct {
	network  string
	addr     string
	protocol protocol.Protocol
	gateway  gateway.Gateway
	logger   *log.Logger
}

// New creates a new server
func New(network, addr string, p protocol.Protocol, gw gateway.Gateway, logger *log.Logger) *Server {
	return &Server{
		network,
		addr,
		p,
		gw,
		logger,
	}
}

// Start starts a server and spawns a goroutine when a new connection is accepted
func (s Server) Start() error {
	ln, err := net.Listen(s.network, s.addr)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			s.logger.Printf("server failed accepting a conn: %s", err.Error())
		} else {
			go s.HandleConn(conn)
		}
	}
}

// HandleConn handles a net.Conn
func (s Server) HandleConn(conn net.Conn) {
	defer func() {
		s.logger.Println("server got a client connection closed")
		conn.Close()
	}()
	s.logger.Println("server got a new client connection")
	r := bufio.NewReader(conn)
	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			return
		}
		if err != nil {
			s.logger.Printf("server failed reading a line: %s", err)
			return
		}
		if keys, err := s.protocol.Parse(line); err != nil {
			s.logger.Printf("server failed parsing a line: %s", err)
		} else {
			for _, k := range keys {
				if key, v, _ := s.gateway.Get(k); v != nil {
					s.protocol.Reply(conn, key, v)
				}
			}
		}
		s.protocol.Finish(conn)
	}
}
