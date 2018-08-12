package server

import (
	"bufio"
	"io"
	"log"
	"net"
)

type OnReadCallbackFunc func(net.Conn, []byte, *log.Logger)

// Server represents a server
type Server struct {
	network string
	addr    string
	logger  *log.Logger
}

// New creates a new server
func New(network, addr string, logger *log.Logger) *Server {
	return &Server{network, addr, logger}
}

// Start starts a server and spawns a goroutine when a new connection is accepted
func (s *Server) Start(callback OnReadCallbackFunc) error {
	ln, err := net.Listen(s.network, s.addr)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			s.logger.Printf("server failed accepting a conn: %s", err.Error())
		} else {
			go s.HandleConn(conn, callback)
		}
	}
}

// HandleConn handles a net.Conn
func (s *Server) HandleConn(conn net.Conn, callback OnReadCallbackFunc) {
	defer conn.Close()
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

		callback(conn, line, s.logger)
	}
}
