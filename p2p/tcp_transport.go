package p2p

import (
	// "bytes"
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

type TCPTransportOps struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOps
	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer

	// listenAddress string
	// listener      net.Listener

	// // handsake funtion ->
	// // handsake Handshaker
	// // handshakeFunc HandshakeFunc
	// shakeHands HandshakeFunc
	// decoder Decoder

	// // common practice is to use import the mu which i want to protact
	// mu    sync.RWMutex // It helps to proctact the shared valible which multiple goroutine is trying to proctact
	// peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOps) *TCPTransport {
	return &TCPTransport{
		TCPTransportOps: opts,
	}
}

// one way to create the tcptranport
// func NewTCPTransport(listenAdrr string) *TCPTransport {
// 	return &TCPTransport{
// 		shakeHands: NOPHandshakeFunc,
// 		// handshakeFunc: NOPHandshakeFunc,
// 		listenAddress: listenAdrr,
// 	}
// }

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
	t.listener, err = net.Listen("tcp", t.ListenAddr)
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
		fmt.Printf("new incoming connection %+v\n", conn)

		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	// creating a peer
	peer := NewTCPPeer(conn, true)

	// if err:=t.shakeHands(peer);err !=nil{
	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error : %s\n", err)
		return
	}

	// lenDecodeError :=0
	// buf := new(bytes.Buffer)
	// Read Loop
	// msg:=&Temp{}
	msg := &Message{}
	// buf := new(bytes.Buffer)/
	// buf := make([]byte, 2000)

	for {
		// n, err := conn.Read(msg)
		// if err != nil {
		// 	fmt.Printf("TCP error: %s\n", err)
		// }
		// msg :=buf[:n]
		if err := t.Decoder.Decode(conn, msg); err != nil {
			// spam protection
			// lenDecodeError++
			// if lenDecodeError == 5{

			// }

			fmt.Printf("TCP  error : %s\n", err)
			continue
		}
		msg.From = conn.RemoteAddr()
		fmt.Printf("Msg : %+s\n", msg)
		fmt.Printf("Message : %+v\n", msg)

		// fmt.Printf("message: %+v\n", buf[:n])
	}

	// fmt.Printf("new incoming connection %+v\n", conn)
}
