package main

import (
	"fmt"
	"log"

	"github.com/unsortedbytes/distributed_file_storage/p2p"
)

func main() {
	// we'r good
	fmt.Println("We Gucci!")

	tcpOpts := p2p.TCPTransportOps{
		ListenAddr:    ":3201",
		Decoder:       p2p.DefaultDecoder{},
		HandshakeFunc: p2p.NOPHandshakeFunc,
	}
	tr := p2p.NewTCPTransport(tcpOpts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}

}
