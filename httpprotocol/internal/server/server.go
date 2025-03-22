package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"

	"github.com/armadi1809/reinventing_the_wheel/httpprotocol/internal/request"
	"github.com/armadi1809/reinventing_the_wheel/httpprotocol/internal/response"
)

type Server struct {
	closed   atomic.Bool
	listener net.Listener
	handler  Handler
}

type HandleError struct {
	Message string
	response.StatusCode
}

type Handler func(w io.Writer, req *request.Request) *HandleError

func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	ser := &Server{
		listener: listener,
		handler:  handler,
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
	handlerBuffer := bytes.NewBuffer(nil)
	request, err := request.RequestFromReader(conn)
	if err != nil {
		herr := &HandleError{
			Message:    err.Error(),
			StatusCode: response.StatusCodeBadRequest,
		}
		herr.write(conn)
		return
	}
	handleErr := s.handler(handlerBuffer, request)
	if handleErr != nil {
		handleErr.write(conn)
		return
	}

	response.WriteStatusLine(conn, response.StatusCodeOk)
	headers := response.GetDefaultHeaders(handlerBuffer.Len())
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		fmt.Printf("error writing headers: %v\n", err)
		return
	}
	conn.Write(handlerBuffer.Bytes())
}

func (he HandleError) write(w io.Writer) {
	status := he.StatusCode
	message := he.Message

	response.WriteStatusLine(w, status)
	headers := response.GetDefaultHeaders(len(message))
	response.WriteHeaders(w, headers)
	w.Write([]byte(message))
}
