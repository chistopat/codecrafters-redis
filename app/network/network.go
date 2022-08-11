package network

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

type Handler = func(in *bufio.Reader, out *bufio.Writer)

type Server struct {
	address  string
	Listener net.Listener
	handler  Handler
}

func NewNetworkServer(address string, port int, handler Handler) *Server {
	addr := fmt.Sprintf("%s:%d", address, port)
	return &Server{
		address: addr,
		handler: handler,
	}
}

func (s *Server) Run() error {
	defer func(s *Server) {
		err := s.Close()
		if err != nil {
			panic("can't close listener!")
		}
	}(s)

	l, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("failed to start listener %s", s.address)
	}
	s.Listener = l
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			return errors.New("could not accept connection")
		}
		if conn == nil {
			return errors.New("could not create connection")
		}
		go s.Handle(conn)
	}
}

func (s *Server) Close() (err error) {
	return s.Listener.Close()
}

func (s *Server) Handle(conn net.Conn) {
	s.handler(bufio.NewReader(conn), bufio.NewWriter(conn))
}
