package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	closed   atomic.Bool
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	ser := &Server{
		listener: listener,
	}
	go ser.listen()
	return ser, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("error occurred when Acceptin the connection %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	message := "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!\n"
	conn.Write([]byte(message))
}
