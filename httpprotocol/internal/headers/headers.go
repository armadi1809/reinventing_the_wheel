package headers

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"unicode"
)

type Headers map[string]string

const crlf = "\r\n"

var validHeaderChars = []rune{'!', '#', '$', '%', '&', '\'', '*', '+', '_', '.', '^', '-', '`', '|', '~'}

func NewHeaders() Headers {
	return make(map[string]string)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	ind := bytes.Index(data, []byte(crlf))
	if ind == -1 {
		return 0, false, nil
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
	if !isHeaderKeyValid(key) {
		return 0, false, fmt.Errorf("invalid header key: %s", key)
	}
	h.set(key, value)
	return ind + 2, false, nil
}

func (h Headers) set(key, val string) {
	lowerKey := strings.ToLower(key)
	if value, ok := h[lowerKey]; ok {
		h[lowerKey] = value + ", " + val
	} else {
		h[lowerKey] = val
	}
}

func isHeaderKeyValid(val string) bool {
	if len(val) == 0 {
		return false
	}
	for _, c := range val {
		if !(unicode.IsDigit(c) || unicode.IsLetter(c) || slices.Contains(validHeaderChars, c)) {
			return false
		}
	}
	return true
}
