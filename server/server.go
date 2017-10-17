package server

import (
	"bufio"
	"log"
	"net"
	"sync"

	"github.com/yowcow/go-romdb/protocol"
	"github.com/yowcow/go-romdb/store"
)

type Server struct {
	proto string
	addr  string

	protocol protocol.Protocol
	store    store.Store

	quit   chan bool
	wg     *sync.WaitGroup
	logger *log.Logger
}

func New(proto, addr string, protocol protocol.Protocol, store store.Store, logger *log.Logger) *Server {
	quit := make(chan bool)
	wg := &sync.WaitGroup{}
	return &Server{
		proto,
		addr,
		protocol,
		store,
		quit,
		wg,
		logger,
	}
}

func (s Server) Start() error {
	defer s.wg.Done()

	ln, err := net.Listen(s.proto, s.addr)
	defer ln.Close()

	if err != nil {
		s.logger.Print(err)
		return err
	}

	nc := make(chan net.Conn)
	go func(l net.Listener, n chan net.Conn) {
		for {
			conn, err := l.Accept()

			if err != nil {
				s.logger.Print("-> failed accepting a new conn: ", err)
			} else {
				n <- conn
			}
		}
	}(ln, nc)

	s.wg.Add(1)
	s.logger.Print("server listening on port: ", s.addr)

	for {
		select {
		case conn := <-nc:
			s.logger.Print("-> accepted a new conn")
			go s.handleConn(conn)
		case <-s.quit:
			s.logger.Print("-> server shutting down")
			s.store.Shutdown()
			return nil
		}
	}
}

func (s Server) Shutdown() error {
	s.quit <- true
	close(s.quit)
	s.wg.Wait()
	return nil
}

func (s Server) handleConn(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)

	for {
		line, _, err := r.ReadLine()

		if err != nil {
			s.logger.Print("-> read error:", err)
			break
		}

		if keys, err := s.protocol.Parse(line); err != nil {
			s.logger.Print("-> parse error: ", err)
		} else {
			for _, k := range keys {
				if v, _ := s.store.Get(k); v != nil {
					s.protocol.Reply(conn, k, v)
				}
			}
		}

		s.protocol.Finish(conn)
	}
}
