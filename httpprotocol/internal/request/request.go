package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	requestLine := strings.Split(string(request), "\r\n")[0]
	requestLineItems := strings.Split(requestLine, " ")
	if len(requestLineItems) != 3 {
		return nil, fmt.Errorf("invalid request line, contains more than 3 parts")
	}
	method := requestLineItems[0]
	if strings.ToUpper(method) != method {
		return nil, fmt.Errorf("invalid http method: %s", method)
	}
	protocol := requestLineItems[2]
	if protocol != "HTTP/1.1" {
		return nil, fmt.Errorf("invalid protocol: %s", protocol)
	}
	target := requestLineItems[1]
	requestLineStrcut := RequestLine{
		HttpVersion:   strings.Split(protocol, "/")[1],
		RequestTarget: target,
		Method:        method,
	}
	return &Request{
		RequestLine: requestLineStrcut,
	}, nil
}
