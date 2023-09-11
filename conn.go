package ws

import (
	"github.com/fengyeall111/ws/pool"
	"github.com/pawelgaczynski/gain"
)

type WsConn struct {
	gain.Conn
	pool.BytesBufferPool
}
