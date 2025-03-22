package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/armadi1809/reinventing_the_wheel/httpprotocol/internal/headers"
	"github.com/armadi1809/reinventing_the_wheel/httpprotocol/internal/request"
	"github.com/armadi1809/reinventing_the_wheel/httpprotocol/internal/response"
	"github.com/armadi1809/reinventing_the_wheel/httpprotocol/internal/server"
)

const port = 42069

const htmlTemplate = `<html>
  <head>
    <title>%s</title>
  </head>
  <body>
    <h1>%s</h1>
    <p>%s</p>
  </body>
</html>`

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

func handler(w *response.Writer, req *request.Request) {
	var headers headers.Headers
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		body := "Your request honestly kinda sucked."
		title := "400 Bad Request"
		header := "Bad Request"
		message := fmt.Sprintf(htmlTemplate, title, header, body)
		w.WriteStatusLine(response.StatusCodeBadRequest)
		headers = response.GetDefaultHeaders(len(message))
		w.WriteHeaders(headers)
		w.WriteBody([]byte(message))
	case "/myproblem":
		body := "Okay, you know what? This one is on me."
		title := "500 Internal Server Error"
		header := "Internal Server Error"
		message := fmt.Sprintf(htmlTemplate, title, header, body)
		w.WriteStatusLine(response.StatusCodeInternalError)
		headers = response.GetDefaultHeaders(len(message))
		w.WriteHeaders(headers)
		w.WriteBody([]byte(message))
	default:
		body := "Your request was an absolute banger."
		title := "200 OK"
		header := "Success!"
		message := fmt.Sprintf(htmlTemplate, title, header, body)
		w.WriteStatusLine(response.StatusCodeOk)
		headers = response.GetDefaultHeaders(len(message))
		w.WriteHeaders(headers)
		w.WriteBody([]byte(message))
	}
}
