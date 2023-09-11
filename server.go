package ws

import (
	"bufio"
	"github.com/fengyeall111/ws/pool"
	"github.com/pawelgaczynski/gain"
	"log/slog"
)

type WebSocketServer interface {
}

type websocketServer struct {
	gain.Server
	slog.Logger
	readerPool pool.FixedPool[*bufio.Reader]
	writerPool pool.FixedPool[*bufio.Writer]
	Upgrader
}

func (w websocketServer) OnStart(server gain.Server) {
	w.Server = server
}

func (w websocketServer) OnAccept(c gain.Conn) {

}

func (w websocketServer) OnRead(c gain.Conn, n int) {
	conn, ok := c.Context().(*WsConn)
	if !ok {
        w.upgrade(c)
	}

	buf, err := c.Next(n)
	if err != nil {
		_ = c.Close()
	}
	// 尝试从buf中解析出完整的帧
}

func(w websocketServer) upgrade(c gain.Conn) (*WsConn,error) {
	var reader = bufio.NewReaderSize(c, c.InboundBuffered())
	rl, err := readLine(reader)
	if err != nil {
		return nil, err
	}
	req, err := httpParseRequestLine(rl)

	if req.major != 1 || req.minor < 1 {
		return nil,
	}

}

func (w websocketServer) OnWrite(c gain.Conn, n int) {
	//TODO implement me
	panic("implement me")
}

func (w websocketServer) OnClose(c gain.Conn, err error) {
	//TODO implement me
	panic("implement me")
}

var _ WebSocketServer = (*websocketServer)(nil)
var _ gain.EventHandler = (*websocketServer)(nil)
