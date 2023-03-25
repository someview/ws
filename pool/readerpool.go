package pool

import (
	"arena"
	"bufio"
	"io"
	"sync"
)

type ReaderPool interface {
	Get(r io.Reader) *bufio.Reader
	Put(br *bufio.Reader)
	Destory() // you can not use Get or Put After call  this method called
}

type readerPool struct {
	*arena.Arena
	*sync.Pool
	size int
}

func (rp *readerPool) Get(r io.Reader) *bufio.Reader {
	br := rp.Pool.Get().(*bufio.Reader)
	br.Reset(r)
	return br
}

func (rp *readerPool) Put(br *bufio.Reader) {
	br.Reset(nil)
	rp.Pool.Put(br)
}

func (rp *readerPool) Destory() {
	rp.Free()
}
func NewReaderPool(bufSize int) ReaderPool {
	a := arena.NewArena()
	p := arena.New[sync.Pool](a)
	p.New = func() any {
		return bufio.NewReaderSize(defaultReadWriter, bufSize)
	}
	return &readerPool{
		Arena: a,
		Pool:  p,
	}
}
