package p2p

import (
	// "bytes"
	"errors"
	"fmt"
	"log"
	"net"
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

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

// Remote Addr implements the Peerinterface and will return  the
// remote addresss of the underlying connection.
func (p *TCPPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

// Close implements the peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder

	OnPeer func(Peer) error
	// OnPeer func(*TCPPeer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC

	// server is resposible to maintaining the peer
	// mu    sync.RWMutex
	// peers map[net.Addr]Peer

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

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

// Consume implements the Transport interface, which will return read-only channedl
// for reading  the incoming message  recived  from the another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

// Close  implements the Transport interface.
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial implements the Transport interface.
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		return err
	}

	// fmt.Println(conn)

	go t.handleConn(conn, true)

	return nil
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

	// log
	log.Printf("TCP tranport listing on the port: %s\n", t.ListenAddr)

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		// creating a peer
		// peer :=NewTCPPeer(conn, true)
		fmt.Printf("new incoming connection %+v\n", conn)

		go t.handleConn(conn, false)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {

	var err error
	defer func() {
		fmt.Printf("dropping peer connection: %s", err)
		conn.Close()
	}()

	// creating a peer
	peer := NewTCPPeer(conn, true)

	// if err:=t.shakeHands(peer);err !=nil{
	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error : %s\n", err)
		return
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
	}

	// lenDecodeError :=0
	// buf := new(bytes.Buffer)
	// Read Loop
	// msg:=&Temp{}
	// msg := &Message{}
	// buf := new(bytes.Buffer)/
	// buf := make([]byte, 2000)
	// rpc := &RPC{} // its the pointer i don't need it
	rpc := RPC{}

	for {
		// n, err := conn.Read(msg)
		// if err != nil {
		// 	fmt.Printf("TCP error: %s\n", err)
		// }
		// msg :=buf[:n]
		err = t.Decoder.Decode(conn, &rpc)
		// fmt.Println(reflect.TypeOf(err))
		// panic(err)

		// if err == net.ErrClosed {
		// 	return
		// }

		// if err == net.OpError{
		// 	return
		// }

		if err != nil {
			// spam protection
			// lenDecodeError++
			// if lenDecodeError == 5{

			// }

			fmt.Printf("TCP read  error : %s\n", err)
			// continue
			return
		}
		rpc.From = conn.RemoteAddr()
		fmt.Printf("Msg : %+s\n", rpc)
		fmt.Printf("Message : %+v\n", rpc)
		t.rpcch <- rpc

		// fmt.Printf("message: %+v\n", buf[:n])
	}

	// fmt.Printf("new incoming connection %+v\n", conn)
}
