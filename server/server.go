package server

import (
	"bufio"
	"log"
	"net"
	"sync"

	"github.com/yowcow/goromdb/protocol"
	"github.com/yowcow/goromdb/store"
)

// Server represents a server
type Server struct {
	proto    string
	addr     string
	protocol protocol.Protocol
	store    store.Store
	logger   *log.Logger
	quit     chan bool
	wg       *sync.WaitGroup
}

// New creates a new server
func New(proto, addr string, protocol protocol.Protocol, store store.Store, logger *log.Logger) *Server {
	quit := make(chan bool)
	wg := &sync.WaitGroup{}
	return &Server{
		proto,
		addr,
		protocol,
		store,
		logger,
		quit,
		wg,
	}
}

// Start starts a server and spawns a goroutine when a new connection is accepted
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
				s.logger.Print("server failed accepting a new conn: ", err)
			} else {
				n <- conn
			}
		}
	}(ln, nc)

	s.wg.Add(1)
	s.logger.Print("server started listening to addr: ", s.addr)

	for {
		select {
		case conn := <-nc:
			s.logger.Print("server accepted a new conn")
			go s.handleConn(conn)
		case <-s.quit:
			s.store.Shutdown()
			s.logger.Print("server finished")
			return nil
		}
	}
}

// Shutdown terminates a server
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
			s.logger.Print("server failed reading a line:", err)
			break
		}
		if keys, err := s.protocol.Parse(line); err != nil {
			s.logger.Print("server failed parsing a line: ", err)
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
