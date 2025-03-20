package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

const crlf = "\r\n"

func NewHeaders() Headers {
	return make(map[string]string)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	ind := bytes.Index(data, []byte(crlf))
	if ind == -1 {
		return 0, false, fmt.Errorf("no carriage returned were found, misformatted headers")
	}
	if ind == 0 {
		return 2, true, nil
	}
	dataToParse := data[:ind]
	columnInd := bytes.Index(dataToParse, []byte(":"))
	key := string(dataToParse[:columnInd])
	value := string(dataToParse[columnInd+1:])
	if strings.HasSuffix(key, " ") {
		return 0, false, fmt.Errorf("end of header key can't contain whitespaces")
	}
	key = strings.TrimSpace(key)
	value = strings.TrimSpace(value)
	h[key] = value
	return ind + 2, false, nil
}
