package main

import (
	"fmt"
	"log"

	"github.com/unsortedbytes/distributed_file_storage/p2p"
)

func main() {
	// we'r good
	fmt.Println("We Gucci!")

	tr:=p2p.NewTCPTransport(":3201")
	if err :=tr.ListenAndAccept(); err!=nil{
		log.Fatal(err)
	}
	select{}

}
