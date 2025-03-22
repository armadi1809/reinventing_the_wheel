package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/armadi1809/reinventing_the_wheel/httpprotocol/internal/headers"
)

type StatusCode int

type Writer struct {
	Writer io.Writer
}

const (
	StatusCodeOk            StatusCode = 200
	StatusCodeBadRequest    StatusCode = 400
	StatusCodeInternalError StatusCode = 500
)

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	_, err := WriteStatusLine(w.Writer, statusCode)
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	return WriteHeaders(w.Writer, headers)
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.Writer.Write(p)
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) (int, error) {
	switch statusCode {
	case StatusCodeOk:
		return w.Write([]byte("HTTP/1.1 200 OK\r\n"))
	case StatusCodeBadRequest:
		return w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
	case StatusCodeInternalError:
		return w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
	}
	return w.Write([]byte("\r\n"))
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("content-length", strconv.Itoa(contentLen))
	headers.Set("connection", "close")
	headers.Set("content-type", "text/html")
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		_, err := w.Write(fmt.Appendf(nil, "%s: %s\r\n", key, val))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
