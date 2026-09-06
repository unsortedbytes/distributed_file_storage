package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListenAddr:    ":3201",
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}
	listenAddr := ":3201"
	tr := NewTCPTransport(opts)

	assert.Equal(t, tr.ListenAddr, listenAddr)

	// Server
	// tr.Start()

	assert.Nil(t, tr.ListenAndAccept())
}
