package main

import (
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/armadi1809/reinventing_the_wheel/httpprotocol/internal/request"
	"github.com/armadi1809/reinventing_the_wheel/httpprotocol/internal/response"
	"github.com/armadi1809/reinventing_the_wheel/httpprotocol/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func handler(w io.Writer, req *request.Request) *server.HandleError {
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		return &server.HandleError{StatusCode: response.StatusCodeBadRequest, Message: "Your problem is not my problem\n"}
	case "/myproblem":
		return &server.HandleError{StatusCode: response.StatusCodeInternalError, Message: "Woopsie, my bad\n"}
	default:
		w.Write([]byte("All good, frfr\n"))
		return nil
	}
}
