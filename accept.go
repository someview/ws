package ws

import "time"

// AcceptOptions represents Accept's options.
type AcceptOptions struct {
	// HandshakeTimeout specifies the duration for the handshake to complete.
	HandshakeTimeout time.Duration
	// Subprotocols lists the WebSocket subprotocols that Accept will negotiate with the client.
	// The empty subprotocol will always be negotiated as per RFC 6455. If you would like to
	// reject it, close the connection when c.Subprotocol() == "".
	SubProtocols []string

	// InsecureSkipVerify is used to disable Accept's origin verification behaviour.
	//
	// You probably want to use OriginPatterns instead.
	InsecureSkipVerify bool

	// OriginPatterns lists the host patterns for authorized origins.
	// The request host is always authorized.
	// Use this to enable cross origin WebSockets.
	//
	// i.e javascript running on example.com wants to access a WebSocket server at chat.example.com.
	// In such a case, example.com is the origin and chat.example.com is the request host.
	// One would set this field to []string{"example.com"} to authorize example.com to connect.
	//
	// Each pattern is matched case insensitively against the request origin host
	// with filepath.Match.
	// See https://golang.org/pkg/path/filepath/#Match
	//
	// Please ensure you understand the ramifications of enabling this.
	// If used incorrectly your WebSocket server will be open to CSRF attacks.
	//
	// Do not use * as a pattern to allow any origin, prefer to use InsecureSkipVerify instead
	// to bring attention to the danger of such a setting.
	OriginPatterns []string

	// CompressionMode controls the compression mode.
	// Defaults to CompressionNoContextTakeover.
	//
	// See docs on CompressionMode for details.
	CompressionMode CompressionMode

	// CompressionThreshold controls the minimum size of a message before compression is applied.
	//
	// Defaults to 512 bytes for CompressionNoContextTakeover and 128 bytes
	// for CompressionContextTakeover.
	CompressionThreshold int
}
