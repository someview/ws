package ws

import (
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/textproto"
	"strings"
)

func headerTokens(h http.Header, key string) []string {
	key = textproto.CanonicalMIMEHeaderKey(key)
	var tokens []string
	for _, v := range h[key] {
		v = strings.TrimSpace(v)
		for _, t := range strings.Split(v, ",") {
			t = strings.TrimSpace(t)
			tokens = append(tokens, t)
		}
	}
	return tokens
}

func headerContainsTokenIgnoreCase(h http.Header, key, token string) bool {
	for _, t := range headerTokens(h, key) {
		if strings.EqualFold(t, token) {
			return true
		}
	}
	return false
}

func isHTTP11OrGreater(str string) bool {
	var protoMajor, protoMinor int
	versions := strings.Split(str, "/")
	if len(versions) < 0 || len(versions) >= 3 {

	}
	return req.ProtoMajor > 1 || (req.ProtoMajor == 1 && req.ProtoMinor >= 1)
}

func verifyClientRequest(ctx fasthttp.RequestCtx) (errCode int, _ error) {
	r := ctx.Request.Header.Protocol()
	r.RequestURI()
	if !ctx.IsGet() || !ctx.Request.Header.IsHTTP11() {
		return http.StatusBadRequest, fmt.Errorf("WebSocket protocol violation: handshake request must be at least HTTP/1.1 or Get")
	}
	var req *http.Request
	if !req.ProtoAtLeast(1, 1) {
		return http.StatusUpgradeRequired, fmt.Errorf("WebSocket protocol violation: handshake request must be at least HTTP/1.1: %q", r.Proto)
	}

	if !headerContainsTokenIgnoreCase(r.Header, "Connection", "Upgrade") {
		w.Header().Set("Connection", "Upgrade")
		w.Header().Set("Upgrade", "websocket")
		return http.StatusUpgradeRequired, fmt.Errorf("WebSocket protocol violation: Connection header %q does not contain Upgrade", r.Header.Get("Connection"))
	}

	if !headerContainsTokenIgnoreCase(r.Header, "Upgrade", "websocket") {
		w.Header().Set("Connection", "Upgrade")
		w.Header().Set("Upgrade", "websocket")
		return http.StatusUpgradeRequired, fmt.Errorf("WebSocket protocol violation: Upgrade header %q does not contain websocket", r.Header.Get("Upgrade"))
	}

	if r.Method != "GET" {
		return http.StatusMethodNotAllowed, fmt.Errorf("WebSocket protocol violation: handshake request method is not GET but %q", r.Method)
	}

	if r.Header.Get("Sec-WebSocket-Version") != "13" {
		w.Header().Set("Sec-WebSocket-Version", "13")
		return http.StatusBadRequest, fmt.Errorf("unsupported WebSocket protocol version (only 13 is supported): %q", r.Header.Get("Sec-WebSocket-Version"))
	}

	if r.Header.Get("Sec-WebSocket-Key") == "" {
		return http.StatusBadRequest, errors.New("WebSocket protocol violation: missing Sec-WebSocket-Key")
	}

	return 0, nil
}
