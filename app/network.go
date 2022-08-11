package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
)

type Handler = func(in *bufio.Reader, out *bufio.Writer)

type NetworkServer struct {
	address  string
	Listener net.Listener
	handler  Handler
}

func NewNetworkServer(address string, port int, handler Handler) (*NetworkServer, error) {
	addr := fmt.Sprintf("%s:%d", address, port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to bind to port %d", port)
	}

	return &NetworkServer{
		address:  addr,
		Listener: l,
		handler:  handler,
	}, nil
}

func (s *NetworkServer) Run() error {
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

func (s *NetworkServer) Close() (err error) {
	return s.Listener.Close()
}

func (s *NetworkServer) Handle(conn net.Conn) {
	defer func(s *NetworkServer) {
		err := s.Close()
		if err != nil {
			panic("can't close listener!")
		}
	}(s)
	s.handler(bufio.NewReader(conn), bufio.NewWriter(conn))
}
