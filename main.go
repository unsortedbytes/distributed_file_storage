package main

import (
	"fmt"
	"log"

	"github.com/unsortedbytes/distributed_file_storage/p2p"
)

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	fmt.Println("doing some logic with the peer outside of TCPTransport")
	return nil
	// return fmt.Errorf("failed the onpeer func")
}

func main() {
	// we'r good
	fmt.Println("We Gucci!")

	tcpOpts := p2p.TCPTransportOps{
		ListenAddr:    ":3201",
		Decoder:       p2p.DefaultDecoder{},
		HandshakeFunc: p2p.NOPHandshakeFunc,
		// OnPeer:        func(p2p.Peer) error { return fmt.Errorf("failed the onpeer func") },
		OnPeer: OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Println(msg)
			fmt.Printf("this msg is %+v\n", msg)
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}

}
