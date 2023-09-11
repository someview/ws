package pool

import (
	"bytes"
)

type BytesBufferPool interface {
	Get(size int) *bytes.Buffer
	Put(size int) *bytes.Buffer
}

type FixedPool[T any] interface {
	Get(size int) T
	Put(T)
}