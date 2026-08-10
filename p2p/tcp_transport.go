package p2p

import "net"

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	peers map[string]


}