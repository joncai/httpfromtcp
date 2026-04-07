package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func validateFieldName(name []byte) error {
	for _, c := range name {
		if !isTokenChar(c) {
			return fmt.Errorf("invalid header field name")
		}
	}
	return nil
}

// isTokenChar reports whether b is a tchar (RFC 7230).
func isTokenChar(b byte) bool {
	return b > 0 && b < 128 && (isAlphaNum(b) || strings.ContainsRune("!#$%&'*+-.^_`|~", rune(b)))
}

func isAlphaNum(b byte) bool {
	return 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || '0' <= b && b <= '9'
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0, false, nil
	}

	rawLine := data[:idx]
	n = idx + 2

	line := bytes.TrimSpace(rawLine)
	if len(line) == 0 {
		return n, true, nil
	}

	colon := bytes.IndexByte(line, ':')
	if colon <= 0 {
		return 0, false, fmt.Errorf("invalid header line")
	}

	if bytes.ContainsAny(line[:colon], " \t") {
		return 0, false, fmt.Errorf("invalid header line")
	}

	if err := validateFieldName(line[:colon]); err != nil {
		return 0, false, err
	}

	name := strings.ToLower(string(line[:colon]))
	value := string(bytes.TrimSpace(line[colon+1:]))
	if existing, ok := h[name]; ok && existing != "" {
		h[name] = existing + ", " + value
	} else {
		h[name] = value
	}

	done = len(data) >= n+2 && data[n] == '\r' && data[n+1] == '\n'
	return n, done, nil
}
