package p2p

import "net"

// Message holds any arbitrary data that is beign sent  over the
// each transport b/w two nodes in the network.
type Message struct {
	From    net.Addr
	Payload []byte
}
