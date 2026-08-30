package p2p

// Peer is an interface that represents the remote node. -> storage device or computer -> node
type Peer interface {
}

// Transport is anything that handles the communication between the nodes in the network.
// This can be of the form(TCP, UDP, websockets,...) -> type of protocals
type Transport interface {
	ListenAndAccept() error
}
