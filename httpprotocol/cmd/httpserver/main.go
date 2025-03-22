package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

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
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyHandler(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/video" {
		handlerVids(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	handler200(w, req)
	return
}

func handlerVids(w *response.Writer, req *request.Request) {
	videoData, err := os.ReadFile("./assets/vim.mp4")
	if err != nil {
		fmt.Printf("error accessign the video %v", err)
		handler500(w, req)
	}
	w.WriteStatusLine(response.StatusCodeOk)
	headers := response.GetDefaultHeaders(len(videoData))
	headers.Override("content-type", "video/mp4")
	w.WriteHeaders(headers)
	w.WriteBody(videoData)
}

func handler500(w *response.Writer, _ *request.Request) {
	body := "Okay, you know what? This one is on me."
	title := "500 Internal Server Error"
	header := "Internal Server Error"
	message := fmt.Sprintf(htmlTemplate, title, header, body)
	w.WriteStatusLine(response.StatusCodeInternalError)
	headers := response.GetDefaultHeaders(len(message))
	w.WriteHeaders(headers)
	w.WriteBody([]byte(message))
}

func handler200(w *response.Writer, _ *request.Request) {
	body := "Your request was an absolute banger."
	title := "200 OK"
	header := "Success!"
	message := fmt.Sprintf(htmlTemplate, title, header, body)
	w.WriteStatusLine(response.StatusCodeOk)
	headers := response.GetDefaultHeaders(len(message))
	w.WriteHeaders(headers)
	w.WriteBody([]byte(message))
}

func handler400(w *response.Writer, _ *request.Request) {
	body := "Your request honestly kinda sucked."
	title := "400 Bad Request"
	header := "Bad Request"
	message := fmt.Sprintf(htmlTemplate, title, header, body)
	w.WriteStatusLine(response.StatusCodeBadRequest)
	headers := response.GetDefaultHeaders(len(message))
	w.WriteHeaders(headers)
	w.WriteBody([]byte(message))
}

func proxyHandler(w *response.Writer, req *request.Request) {
	target := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	url := "https://httpbin.org/" + target
	fmt.Println("Proxying to", url)
	resp, err := http.Get(url)
	if err != nil {
		handler500(w, req)
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(response.StatusCodeOk)
	h := response.GetDefaultHeaders(0)
	h.Override("Transfer-Encoding", "chunked")
	h.Remove("Content-Length")
	w.WriteHeaders(h)

	const maxChunkSize = 1024
	buffer := make([]byte, maxChunkSize)
	for {
		n, err := resp.Body.Read(buffer)
		fmt.Println("Read", n, "bytes")
		if n > 0 {
			_, err = w.WriteChunkedBody(buffer[:n])
			if err != nil {
				fmt.Println("Error writing chunked body:", err)
				break
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading response body:", err)
			break
		}
	}
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Println("Error writing chunked body done:", err)
	}
}
