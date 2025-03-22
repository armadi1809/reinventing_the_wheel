package response

import (
	"fmt"
	"io"
	"strconv"

	"github.com/armadi1809/reinventing_the_wheel/httpprotocol/internal/headers"
)

type StatusCode int

const (
	StatusCodeOk            StatusCode = 200
	StatusCodeBadResponse   StatusCode = 400
	StatusCodeInternalError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) {
	switch statusCode {
	case StatusCodeOk:
		w.Write([]byte("HTTP/1.1 200 OK"))
	case StatusCodeBadResponse:
		w.Write([]byte("HTTP/1.1 400 Bad Request"))
	case StatusCodeInternalError:
		w.Write([]byte("HTTP/1.1 500 Internal Server Error"))
	}
	w.Write([]byte("\r\n"))
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers.Set("content-length", strconv.Itoa(contentLen))
	headers.Set("connection", "close")
	headers.Set("content-type", "text/plain")
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		_, err := w.Write(fmt.Appendf(nil, "%s: %s\r\n", key, val))
		if err != nil {
			return err
		}
	}
	w.Write([]byte("\r\n"))
	return nil
}
