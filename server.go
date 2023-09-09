package ws

type OnConnect func(conn WsConn)

type WsServer struct {
}

type WsOption struct {
	Protocols []string
}

func NewServer(opt WsOption) *WsServer {
	return nil
}

func (s *WsServer) Run() error {
	return nil
}
