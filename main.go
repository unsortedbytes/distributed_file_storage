package main

import (
	"fmt"
	"log"

	"github.com/unsortedbytes/distributed_file_storage/p2p"
)

// "fmt"
// "log"

// "github.com/unsortedbytes/distributed_file_storage/p2p"

// func OnPeer(peer p2p.Peer) error {
// 	peer.Close()
// 	fmt.Println("doing some logic with the peer outside of TCPTransport")
// 	return nil
// 	// return fmt.Errorf("failed the onpeer func")
// }

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTranportOpts := p2p.TCPTransportOpts{
		// ListenAddr:    ":3201",
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},

		// Todo : On Peer func
	}

	tcpTransport := p2p.NewTCPTransport(tcpTranportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:        listenAddr + "_network",
		PathTransformsFunc: CASPathTranformsFunc,
		Transport:          tcpTransport,
		// BootstrapNodes:     []string{":4000"},
		BootstrapNodes: nodes,
	}

	s := NewFileServer(fileServerOpts)
	tcpTransport.OnPeer = s.OnPeer
	return s

}

func main() {
	// we'r good
	// fmt.Println("We Gucci!")

	// tcpOpts := p2p.TCPTransportOps{
	// 	ListenAddr:    ":3201",
	// 	Decoder:       p2p.DefaultDecoder{},
	// 	HandshakeFunc: p2p.NOPHandshakeFunc,
	// 	// OnPeer:        func(p2p.Peer) error { return fmt.Errorf("failed the onpeer func") },
	// 	OnPeer: OnPeer,
	// }
	// tr := p2p.NewTCPTransport(tcpOpts)

	// go func() {
	// 	for {
	// 		msg := <-tr.Consume()
	// 		fmt.Println(msg)
	// 		fmt.Printf("this msg is %+v\n", msg)
	// 	}
	// }()

	// if err := tr.ListenAndAccept(); err != nil {
	// 	log.Fatal(err)
	// }
	// select {}

	fmt.Println("We Gucci")
	// tcpTranportOpts := p2p.TCPTransportOpts{
	// 	ListenAddr:    ":3201",
	// 	HandshakeFunc: p2p.NOPHandshakeFunc,
	// 	Decoder:       p2p.DefaultDecoder{},

	// 	// Todo : On Peer func
	// }

	// tcpTransport := p2p.NewTCPTransport(tcpTranportOpts)

	// fileServerOpts := FileServerOpts{
	// 	StorageRoot:        "3201_network",
	// 	PathTransformsFunc: CASPathTranformsFunc,
	// 	Transport:          tcpTransport,
	// 	BootstrapNodes:     []string{":4000"},
	// }

	// s := NewFileServer(fileServerOpts)

	// go func() {
	// 	time.Sleep(time.Second * 3)
	// 	s.Stop()
	// }()

	// if err := s.Start(); err != nil {
	// 	log.Fatal(err)
	// }

	// // select {}

	s1 := makeServer(":3201", "")
	s2 := makeServer(":4000", ":3201")

	go func() {
		log.Fatal(s1.Start())
	}()

	s2.Start()

}
