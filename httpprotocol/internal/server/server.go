package server

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"

	"github.com/armadi1809/reinventing_the_wheel/httpprotocol/internal/response"
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
	response.WriteStatusLine(conn, response.StatusCodeOk)
	headers := response.GetDefaultHeaders(0)
	err := response.WriteHeaders(conn, headers)
	if err != nil {
		return
	}
}
