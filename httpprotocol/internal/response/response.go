package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/armadi1809/reinventing_the_wheel/httpprotocol/internal/headers"
)

type StatusCode int

type Writer struct {
	writer io.Writer
}

const (
	StatusCodeOk            StatusCode = 200
	StatusCodeBadRequest    StatusCode = 400
	StatusCodeInternalError StatusCode = 500
)

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	_, err := WriteStatusLine(w.writer, statusCode)
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	return WriteHeaders(w.writer, headers)
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	chunkSize := len(p)

	nTotal := 0
	n, err := fmt.Fprintf(w.writer, "%x\r\n", chunkSize)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write(p)
	if err != nil {
		return nTotal, err
	}
	nTotal += n

	n, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return nTotal, err
	}
	nTotal += n
	return nTotal, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := w.writer.Write([]byte("0\r\n\r\n"))
	if err != nil {
		return n, err
	}
	return n, nil
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
