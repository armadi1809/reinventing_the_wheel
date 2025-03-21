package main

import (
	"fmt"
	"net"

	"github.com/armadi1809/reinventing_the_wheel/httpprotocol/internal/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic("error setting up the listener")
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			panic("error accepting the connection")
		}
		fmt.Println("Connection accepted")
		request, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Printf("error parsing request from connection %v", err)
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)
		fmt.Println("Headers: ")
		for key, val := range request.Headers {
			fmt.Println("- " + key + ": " + val)
		}
	}

}
