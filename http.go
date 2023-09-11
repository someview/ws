package ws

import (
	"bytes"
	_ "github.com/fengyeall111/gnet-http"
)


var (
	httpVersion1_0    = []byte("HTTP/1.0")
	httpVersion1_1    = []byte("HTTP/1.1")
	httpVersionPrefix = []byte("HTTP/")
)

type httpRequestLine struct {
	method, uri  []byte
	major, minor int
}

type httpResponseLine struct {
	major, minor int
	status       int
	reason       []byte
}

type HandShakeError struct {
	Reason string
	Code   int
	Header map[string]string
}


func httpParseRequestLine(line []byte) (req httpRequestLine, err error) {
	var proto []byte
	req.method, req.uri, proto = bsplit3(line, ' ')

	var ok bool
	req.major, req.minor, ok = httpParseVersion(proto)
	if !ok {
		err = ErrMalformedRequest
		return
	}

	return
}

func httpParseVersion(bts []byte) (major, minor int, ok bool) {
	switch {
	case bytes.Equal(bts, httpVersion1_0):
		return 1, 0, true
	case bytes.Equal(bts, httpVersion1_1):
		return 1, 1, true
	case len(bts) < 8:
		return
	case !bytes.Equal(bts[:5], httpVersionPrefix):
		return
	}

	bts = bts[5:]

	dot := bytes.IndexByte(bts, '.')
	if dot == -1 {
		return
	}
	var err error
	major, err = asciiToInt(bts[:dot])
	if err != nil {
		return
	}
	minor, err = asciiToInt(bts[dot+1:])
	if err != nil {
		return
	}

	return major, minor, true
}