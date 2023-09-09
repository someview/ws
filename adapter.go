package ws

import (
	"github.com/pawelgaczynski/gain"
	"log/slog"
	"time"
)

type adapter struct {
	gain.Server
	slog.Logger
}

func (a *adapter) OnStart(server gain.Server) {
	a.Server = server
}

func (a *adapter) OnAccept(c gain.Conn) {
	err := c.SetKeepAlivePeriod(time.Minute * 2)
	if err != nil {

	}
}

func (a *adapter) OnRead(c gain.Conn, n int) {
	c.
}

func (a *adapter) OnWrite(c gain.Conn, n int) {
	//TODO implement me
	panic("implement me")
}

func (a *adapter) OnClose(c gain.Conn, err error) {
	//TODO implement me
	panic("implement me")
}

var _ gain.EventHandler = (*adapter)(nil)
