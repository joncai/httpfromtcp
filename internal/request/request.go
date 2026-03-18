package request

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Parser state for a Request.
const (
	StateInitialized int = iota
	StateDone
)

type Request struct {
	RequestLine RequestLine
	State       int
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

// ParseRequestLine parses the first HTTP request line from data.
// It returns the RequestLine, the number of bytes consumed (including the trailing \r\n), and any error.
// If \r\n is not found, it returns (zero, 0, nil) to indicate more data is needed.
func ParseRequestLine(data []byte) (RequestLine, int, error) {
	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return RequestLine{}, 0, nil
	}
	line := string(data[:idx])
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return RequestLine{}, 0, fmt.Errorf("invalid request line")
	}
	version := strings.TrimPrefix(parts[2], "HTTP/")
	if version != "1.1" {
		return RequestLine{}, 0, fmt.Errorf("invalid request line")
	}
	return RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   version,
	}, idx + 2, nil
}

// parse accepts unparsed bytes from the buffer, updates the parser state and
// RequestLine, and returns the number of bytes consumed and any error.
func (r *Request) parse(data []byte) (int, error) {
	if r.State == StateDone {
		return 0, nil
	}
	requestLine, n, err := ParseRequestLine(data)
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return 0, nil
	}
	r.RequestLine = requestLine
	r.State = StateDone
	return n, nil
}

const initialBufSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{State: StateInitialized}
	buf := make([]byte, 0, initialBufSize)
	var bytesRead, bytesParsed int

	for r.State != StateDone {
		consumed, err := r.parse(buf)
		if err != nil {
			return nil, err
		}
		bytesParsed += consumed
		if consumed > 0 {
			buf = buf[consumed:]
		}
		if r.State == StateDone {
			return r, nil
		}
		// Need more data: read next chunk
		readBuf := make([]byte, 512)
		n, readErr := reader.Read(readBuf)
		bytesRead += n
		if n > 0 {
			buf = append(buf, readBuf[:n]...)
		}
		if readErr != nil {
			if readErr == io.EOF {
				if len(buf) > 0 {
					return nil, fmt.Errorf("no request line found")
				}
				return nil, fmt.Errorf("no request line found")
			}
			return nil, readErr
		}
	}
	return r, nil
}
