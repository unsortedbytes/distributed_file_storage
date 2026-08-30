package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {
	// conn is the underlying connection of the peer
	conn net.Conn

	// if we dial and  retrive a conn => outbound == true
	// if we accept and retrive a conn => outbound == false

	outbound bool // dial
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	// common practice is to use import the mu which i want to protact
	mu    sync.RWMutex // It helps to proctact the shared valible which multiple goroutine is trying to proctact
	peers map[net.Addr]Peer
}

// one way to create the tcptranport
func NewTCPTransport(listenAdrr string) *TCPTransport {
	return &TCPTransport{
		listenAddress: listenAdrr,
	}
}

/*
func NewTCPTransport(listenAdrr string) Transport{
	return &TCPTransport{
		listenAddress: listenAdrr,
	}
}


func Tset(){
	t := NewTCPTransport(":4344").(*TCPTransport)

	t.listener.Accept()
}

*/

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		// creating a peer
		// peer :=NewTCPPeer(conn, true)

		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	// creating a peer
	// peer :=NewTCPPeer(conn, true)
	fmt.Printf("new incoming connection %+v\n", conn)
}
