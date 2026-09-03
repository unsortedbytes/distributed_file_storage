package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":4000"
	tr := NewTCPTransport(ListenAddr)

	assert.Equal(t, tr.listenAddress, listenAddr)

	// Server
	// tr.Start()

	assert.Nil(t, tr.ListenAndAccept())
}
