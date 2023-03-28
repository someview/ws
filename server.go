package ws

import (
	http "github.com/fengyeall111/gnet-http"
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/fasthttp"
)

type WebSocketServer struct { // 添加一个环形队列
	http.HttpServer
	AcceptOptions                                             // 提供通用的功能
	OnBeforeUpgrade func(ctx fasthttp.RequestCtx) (err error) // 自定义的功能可以在这钩子函数里实现
}

type HandShake func(ctx fasthttp.RequestCtx)

// todo 将所有的项目添加到功能
func (wss *WebSocketServer) Run(addr string, ack AcceptOptions, opts ...gnet.Option) error {
	return wss.HttpServer.Run(addr, opts...)
}

func (wss *WebSocketServer) handShake(ctx fasthttp.RequestCtx) {
}

func (wss *WebSocketServer) OnTextMessage() {
}
func (wss *WebSocketServer) OnBinaryMessage() {

}
func (wss *WebSocketServer) OnPingMessage() {
}
func (wss *WebSocketServer) OnPongMessage() {
}

func (wss *WebSocketServer) WriteTextMessageAsync([]byte) {
}
func (wss *WebSocketServer) WriteBinaryMessageAsync([]byte) {
}
func (wss *WebSocketServer) WritePingMessageAsync([]byte) {
}
func (wss *WebSocketServer) WritePongMessageAsync([]byte) {
}
func (wss *WebSocketServer) WriteCloseMessageAsync([]byte) {
}
