package p2p

import "net"

// Peer is an interface that represents the remote node. -> storage device or computer -> node
type Peer interface {
	// Conn() net.Conn
	net.Conn // *********** This peak golang engineering ****************
	Send([]byte) error
	// RemoteAddr() net.Addr
	// Close() error
}

// Transport is anything that handles the communication between the nodes in the network.
// This can be of the form(TCP, UDP, websockets,...) -> type of protocals
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error

	Dial(string) error
	// ListenAddr() string
}
