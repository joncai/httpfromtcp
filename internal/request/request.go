package request

import (
	"bufio"
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

func ParseRequestLine(line string) (RequestLine, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return RequestLine{}, fmt.Errorf("invalid request line")
	}
	version := strings.TrimPrefix(parts[2], "HTTP/")
	if version != "1.1" {
		return RequestLine{}, fmt.Errorf("invalid request line")
	}
	return RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   version,
	}, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	scanner := bufio.NewScanner(reader)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("no request line found")
	}
	line := scanner.Text()
	if line == "" {
		return nil, fmt.Errorf("no request line found")
	}
	requestLine, err := ParseRequestLine(line)
	if err != nil {
		return nil, err
	}
	return &Request{
		RequestLine: requestLine,
	}, nil
}
