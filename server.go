package ws

import (
	"bytes"
	"fmt"
	"github.com/gobwas/httphead"
	"github.com/gobwas/pool/pbufio"
	"io"
	"net/http"
)

// Constants used by ConnUpgrader.
const (
	DefaultServerReadBufferSize  = 4096
	DefaultServerWriteBufferSize = 512
)

// Errors used by both client and server when preparing WebSocket handshake.
var (
	ErrHandshakeBadProtocol = RejectConnectionError(
		RejectionStatus(http.StatusHTTPVersionNotSupported),
		RejectionReason(fmt.Sprintf("handshake error: bad HTTP protocol version")),
	)
	ErrHandshakeBadMethod = RejectConnectionError(
		RejectionStatus(http.StatusMethodNotAllowed),
		RejectionReason(fmt.Sprintf("handshake error: bad HTTP request method")),
	)
	ErrHandshakeBadHost = RejectConnectionError(
		RejectionStatus(http.StatusBadRequest),
		RejectionReason(fmt.Sprintf("handshake error: bad %q header", headerHost)),
	)
	ErrHandshakeBadUpgrade = RejectConnectionError(
		RejectionStatus(http.StatusBadRequest),
		RejectionReason(fmt.Sprintf("handshake error: bad %q header", headerUpgrade)),
	)
	ErrHandshakeBadConnection = RejectConnectionError(
		RejectionStatus(http.StatusBadRequest),
		RejectionReason(fmt.Sprintf("handshake error: bad %q header", headerConnection)),
	)
	ErrHandshakeBadSecAccept = RejectConnectionError(
		RejectionStatus(http.StatusBadRequest),
		RejectionReason(fmt.Sprintf("handshake error: bad %q header", headerSecAccept)),
	)
	ErrHandshakeBadSecKey = RejectConnectionError(
		RejectionStatus(http.StatusBadRequest),
		RejectionReason(fmt.Sprintf("handshake error: bad %q header", headerSecKey)),
	)
	ErrHandshakeBadSecVersion = RejectConnectionError(
		RejectionStatus(http.StatusBadRequest),
		RejectionReason(fmt.Sprintf("handshake error: bad %q header", headerSecVersion)),
	)
)

// ErrMalformedResponse is returned by Dialer to indicate that server response
// can not be parsed.
var ErrMalformedResponse = fmt.Errorf("malformed HTTP response")

// ErrMalformedRequest is returned when HTTP request can not be parsed.
var ErrMalformedRequest = RejectConnectionError(
	RejectionStatus(http.StatusBadRequest),
	RejectionReason("malformed HTTP request"),
)

// ErrHandshakeUpgradeRequired is returned by Upgrader to indicate that
// connection is rejected because given WebSocket version is malformed.
//
// According to RFC6455:
// If this version does not match a version understood by the server, the
// server MUST abort the WebSocket handshake described in this section and
// instead send an appropriate HTTP error code (such as 426 Upgrade Required)
// and a |Sec-WebSocket-Version| header field indicating the version(s) the
// server is capable of understanding.
var ErrHandshakeUpgradeRequired = RejectConnectionError(
	RejectionStatus(http.StatusUpgradeRequired),
	RejectionHeader(HandshakeHeaderString(headerSecVersion+": 13\r\n")),
	RejectionReason(fmt.Sprintf("handshake error: bad %q header", headerSecVersion)),
)

// ErrNotHijacker is an error returned when http.ResponseWriter does not
// implement http.Hijacker interface.
var ErrNotHijacker = RejectConnectionError(
	RejectionStatus(http.StatusInternalServerError),
	RejectionReason("given http.ResponseWriter is not a http.Hijacker"),
)

// DefaultUpgrader is an Upgrader that holds no options and is used by Upgrade
// function.
var DefaultUpgrader Upgrader

// Upgrade is like Upgrader{}.Upgrade().
func Upgrade(conn io.ReadWriter) (Handshake, error) {
	return DefaultUpgrader.Upgrade(conn)
}

// Upgrade upgrades http connection to the websocket connection.
// Upgrader contains options for upgrading connection to websocket.
type Upgrader struct {
	// ReadBufferSize and WriteBufferSize is an I/O buffer sizes.
	// They used to read and write http data while upgrading to WebSocket.
	// Allocated buffers are pooled with sync.Pool to avoid extra allocations.
	//
	// If a size is zero then default value is used.
	//
	// Usually it is useful to set read buffer size bigger than write buffer
	// size because incoming request could contain long header values, such as
	// Cookie. Response, in other way, could be big only if user write multiple
	// custom headers. Usually response takes less than 256 bytes.
	ReadBufferSize, WriteBufferSize int

	// Protocol is a select function that is used to select subprotocol
	// from list requested by client. If this field is set, then the first matched
	// protocol is sent to a client as negotiated.
	//
	// The argument is only valid until the callback returns.
	Protocol func([]byte) bool

	// ProtocolCustrom allow user to parse Sec-WebSocket-Protocol header manually.
	// Note that returned bytes must be valid until Upgrade returns.
	// If ProtocolCustom is set, it used instead of Protocol function.
	ProtocolCustom func([]byte) (string, bool)
	
	// ExtensionCustom allow user to parse Sec-WebSocket-Extensions header
	// manually.
	//
	// If ExtensionCustom() decides to accept received extension, it must
	// append appropriate option to the given slice of httphead.Option.
	// It returns results of append() to the given slice and a flag that
	// reports whether given header value is wellformed or not.
	//
	// Note that ExtensionCustom may be called multiple times and
	// implementations must track uniqueness of accepted extensions manually.
	//
	// Note that returned options should be valid until Upgrade returns.
	// If ExtensionCustom is set, it used instead of Extension function.
	ExtensionCustom func([]byte, []httphead.Option) ([]httphead.Option, bool)

	// Negotiate is the callback that is used to negotiate extensions from
	// the client's offer. If this field is set, then the returned non-zero
	// extensions are sent to the client as accepted extensions in the
	// response.
	//
	// The argument is only valid until the Negotiate callback returns.
	//
	// If returned error is non-nil then connection is rejected and response is
	// sent with appropriate HTTP error code and body set to error message.
	//
	// RejectConnectionError could be used to get more control on response.
	Negotiate func(httphead.Option) (httphead.Option, error)

	// Header is an optional HandshakeHeader instance that could be used to
	// write additional headers to the handshake response.
	//
	// It used instead of any key-value mappings to avoid allocations in user
	// land.
	//
	// Note that if present, it will be written in any result of handshake.
	Header HandshakeHeader

	// OnRequest is a callback that will be called after request line
	// successful parsing.
	//
	// The arguments are only valid until the callback returns.
	//
	// If returned error is non-nil then connection is rejected and response is
	// sent with appropriate HTTP error code and body set to error message.
	//
	// RejectConnectionError could be used to get more control on response.
	OnRequest func(uri []byte) error

	// OnHost is a callback that will be called after "Host" header successful
	// parsing.
	//
	// It is separated from OnHeader callback because the Host header must be
	// present in each request since HTTP/1.1. Thus Host header is non-optional
	// and required for every WebSocket handshake.
	//
	// The arguments are only valid until the callback returns.
	//
	// If returned error is non-nil then connection is rejected and response is
	// sent with appropriate HTTP error code and body set to error message.
	//
	// RejectConnectionError could be used to get more control on response.
	OnHost func(host []byte) error

	// OnHeader is a callback that will be called after successful parsing of
	// header, that is not used during WebSocket handshake procedure. That is,
	// it will be called with non-websocket headers, which could be relevant
	// for application-level logic.
	//
	// The arguments are only valid until the callback returns.
	//
	// If returned error is non-nil then connection is rejected and response is
	// sent with appropriate HTTP error code and body set to error message.
	//
	// RejectConnectionError could be used to get more control on response.
	OnHeader func(key, value []byte) error

	// OnBeforeUpgrade is a callback that will be called before sending
	// successful upgrade response.
	//
	// Setting OnBeforeUpgrade allows user to make final application-level
	// checks and decide whether this connection is allowed to successfully
	// upgrade to WebSocket.
	//
	// It must return non-nil either HandshakeHeader or error and never both.
	//
	// If returned error is non-nil then connection is rejected and response is
	// sent with appropriate HTTP error code and body set to error message.
	//
	// RejectConnectionError could be used to get more control on response.
	OnBeforeUpgrade func() (header HandshakeHeader, err error)
}

// Upgrade zero-copy upgrades connection to WebSocket. It interprets given conn
// as connection with incoming HTTP Upgrade request.
//
// It is a caller responsibility to manage i/o timeouts on conn.
//
// Non-nil error means that request for the WebSocket upgrade is invalid or
// malformed and usually connection should be closed.
// Even when error is non-nil Upgrade will write appropriate response into
// connection in compliance with RFC.
func (u Upgrader) Upgrade(conn io.ReadWriter) (hs Handshake, err error) {
	// headerSeen constants helps to report whether or not some header was seen
	// during reading request bytes.
	const (
		headerSeenHost = 1 << iota
		headerSeenUpgrade
		headerSeenConnection
		headerSeenSecVersion
		headerSeenSecKey

		// headerSeenAll is the value that we expect to receive at the end of
		// headers read/parse loop.
		headerSeenAll = 0 |
			headerSeenHost |
			headerSeenUpgrade |
			headerSeenConnection |
			headerSeenSecVersion |
			headerSeenSecKey
	)

	// Prepare I/O buffers.
	// TODO(gobwas): make it configurable.
	br := pbufio.GetReader(conn,
		nonZero(u.ReadBufferSize, DefaultServerReadBufferSize),
	)
	bw := pbufio.GetWriter(conn,
		nonZero(u.WriteBufferSize, DefaultServerWriteBufferSize),
	)
	defer func() {
		pbufio.PutReader(br)
		pbufio.PutWriter(bw)
	}()

	// Read HTTP request line like "GET /ws HTTP/1.1".
	rl, err := readLine(br)
	if err != nil {
		return
	}
	// Parse request line data like HTTP version, uri and method.
	req, err := httpParseRequestLine(rl)
	if err != nil {
		return
	}

	// Prepare stack-based handshake header list.
	header := handshakeHeader{
		0: u.Header,
	}

	// Parse and check HTTP request.
	// As RFC6455 says:
	//   The client's opening handshake consists of the following parts. If the
	//   server, while reading the handshake, finds that the client did not
	//   send a handshake that matches the description below (note that as per
	//   [RFC2616], the order of the header fields is not important), including
	//   but not limited to any violations of the ABNF grammar specified for
	//   the components of the handshake, the server MUST stop processing the
	//   client's handshake and return an HTTP response with an appropriate
	//   error code (such as 400 Bad Request).
	//
	// See https://tools.ietf.org/html/rfc6455#section-4.2.1

	// An HTTP/1.1 or higher GET request, including a "Request-URI".
	//
	// Even if RFC says "1.1 or higher" without mentioning the part of the
	// version, we apply it only to minor part.
	switch {
	case req.major != 1 || req.minor < 1:
		// Abort processing the whole request because we do not even know how
		// to actually parse it.
		err = ErrHandshakeBadProtocol

	case btsToString(req.method) != http.MethodGet:
		err = ErrHandshakeBadMethod

	default:
		if onRequest := u.OnRequest; onRequest != nil {
			err = onRequest(req.uri)
		}
	}
	// Start headers read/parse loop.
	var (
		// headerSeen reports which header was seen by setting corresponding
		// bit on.
		headerSeen byte

		nonce = make([]byte, nonceSize)
	)
	for err == nil {
		line, e := readLine(br)
		if e != nil {
			return hs, e
		}
		if len(line) == 0 {
			// Blank line, no more lines to read.
			break
		}

		k, v, ok := httpParseHeaderLine(line)
		if !ok {
			err = ErrMalformedRequest
			break
		}

		switch btsToString(k) {
		case headerHostCanonical:
			headerSeen |= headerSeenHost
			if onHost := u.OnHost; onHost != nil {
				err = onHost(v)
			}

		case headerUpgradeCanonical:
			headerSeen |= headerSeenUpgrade
			if !bytes.Equal(v, specHeaderValueUpgrade) && !bytes.EqualFold(v, specHeaderValueUpgrade) {
				err = ErrHandshakeBadUpgrade
			}

		case headerConnectionCanonical:
			headerSeen |= headerSeenConnection
			if !bytes.Equal(v, specHeaderValueConnection) && !btsHasToken(v, specHeaderValueConnectionLower) {
				err = ErrHandshakeBadConnection
			}

		case headerSecVersionCanonical:
			headerSeen |= headerSeenSecVersion
			if !bytes.Equal(v, specHeaderValueSecVersion) {
				err = ErrHandshakeUpgradeRequired
			}

		case headerSecKeyCanonical:
			headerSeen |= headerSeenSecKey
			if len(v) != nonceSize {
				err = ErrHandshakeBadSecKey
			} else {
				copy(nonce[:], v)
			}

		case headerSecProtocolCanonical:
			if custom, check := u.ProtocolCustom, u.Protocol; hs.Protocol == "" && (custom != nil || check != nil) {
				var ok bool
				if custom != nil {
					hs.Protocol, ok = custom(v)
				} else {
					hs.Protocol, ok = btsSelectProtocol(v, check)
				}
				if !ok {
					err = ErrMalformedRequest
				}
			}

		case headerSecExtensionsCanonical:
			if f := u.Negotiate; err == nil && f != nil {
				hs.Extensions, err = negotiateExtensions(v, hs.Extensions, f)
			}
			// DEPRECATED path.
			if custom, check := u.ExtensionCustom, u.Extension; u.Negotiate == nil && (custom != nil || check != nil) {
				var ok bool
				if custom != nil {
					hs.Extensions, ok = custom(v, hs.Extensions)
				} else {
					hs.Extensions, ok = btsSelectExtensions(v, hs.Extensions, check)
				}
				if !ok {
					err = ErrMalformedRequest
				}
			}

		default:
			if onHeader := u.OnHeader; onHeader != nil {
				err = onHeader(k, v)
			}
		}
	}
	switch {
	case err == nil && headerSeen != headerSeenAll:
		switch {
		case headerSeen&headerSeenHost == 0:
			// As RFC2616 says:
			//   A client MUST include a Host header field in all HTTP/1.1
			//   request messages. If the requested URI does not include an
			//   Internet host name for the service being requested, then the
			//   Host header field MUST be given with an empty value. An
			//   HTTP/1.1 proxy MUST ensure that any request message it
			//   forwards does contain an appropriate Host header field that
			//   identifies the service being requested by the proxy. All
			//   Internet-based HTTP/1.1 servers MUST respond with a 400 (Bad
			//   Request) status code to any HTTP/1.1 request message which
			//   lacks a Host header field.
			err = ErrHandshakeBadHost
		case headerSeen&headerSeenUpgrade == 0:
			err = ErrHandshakeBadUpgrade
		case headerSeen&headerSeenConnection == 0:
			err = ErrHandshakeBadConnection
		case headerSeen&headerSeenSecVersion == 0:
			// In case of empty or not present version we do not send 426 status,
			// because it does not meet the ABNF rules of RFC6455:
			//
			// version = DIGIT | (NZDIGIT DIGIT) |
			// ("1" DIGIT DIGIT) | ("2" DIGIT DIGIT)
			// ; Limited to 0-255 range, with no leading zeros
			//
			// That is, if version is really invalid – we sent 426 status as above, if it
			// not present – it is 400.
			err = ErrHandshakeBadSecVersion
		case headerSeen&headerSeenSecKey == 0:
			err = ErrHandshakeBadSecKey
		default:
			panic("unknown headers state")
		}

	case err == nil && u.OnBeforeUpgrade != nil:
		header[1], err = u.OnBeforeUpgrade()
	}
	if err != nil {
		var code int
		if rej, ok := err.(*ConnectionRejectedError); ok {
			code = rej.code
			header[1] = rej.header
		}
		if code == 0 {
			code = http.StatusInternalServerError
		}
		httpWriteResponseError(bw, err, code, header.WriteTo)
		// Do not store Flush() error to not override already existing one.
		_ = bw.Flush()
		return
	}

	httpWriteResponseUpgrade(bw, nonce, hs, header.WriteTo)
	err = bw.Flush()

	return
}

type handshakeHeader [2]HandshakeHeader

func (hs handshakeHeader) WriteTo(w io.Writer) (n int64, err error) {
	for i := 0; i < len(hs) && err == nil; i++ {
		if h := hs[i]; h != nil {
			var m int64
			m, err = h.WriteTo(w)
			n += m
		}
	}
	return n, err
}
