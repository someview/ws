package pool

import (
	"io"
)

type simpleReadWriter struct{}

func (s *simpleReadWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}
func (s *simpleReadWriter) Read(p []byte) (n int, err error) {
	return 0, nil
}
func NewReadWriter() io.ReadWriter {
	return &simpleReadWriter{}
}

var defaultReadWriter = &simpleReadWriter{} // just used for init
