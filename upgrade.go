package ws

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

type handshakeRequest struct {
}

type handshakeResponse struct {
}