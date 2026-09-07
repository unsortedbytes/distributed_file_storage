package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/unsortedbytes/distributed_file_storage/p2p"
)

type FileServerOpts struct {
	// ListenAddr         string
	StorageRoot        string
	PathTransformsFunc PathTransformsFunc
	Transport          p2p.Transport
	// TCPTransportOps p2p.TCPTransportOps
	BootstrapNodes []string
}

type FileServer struct {
	FileServerOpts

	// adding the peer
	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	store  *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTranformsFunc: opts.PathTransformsFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),

		peers: make(map[string]p2p.Peer),
	}
}

type Message struct{
	From string
	Payload any
}

// type Payload struct {
// 	Key  string
// 	Data []byte
// }
type DataMessage struct{
	Key string
	Data []byte
}

// func (s *FileServer) broadcast(p *DataMessage) error {
func (s *FileServer) broadcast(msg *Message) error {
	// buf := new(bytes.Buffer)
	// for _ , peer := range s.peers{
	// 	if err := gob.NewEncoder(buf).Encode(p); err!=nil{
	// 		return err
	// 	}

	// 	peer.Send(buf.Bytes())
	// }

	// for _ , peer := range s.peers{
	// 	if err := gob.NewEncoder(peer).Encode(p); err!=nil{
	// 		return err
	// 	}

	// }

	// return nil

	// peers := []p2p.Peer{}
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	// 1. Store this file to disk
	// 2. broadcast this file to all know peers in the network

	// Making the buf and tree here before write
	buf := new(bytes.Buffer)
	tee := io.TeeReader(r, buf)

	if err := s.store.Write(key, tee); err != nil {
		return err
	}

	// the reader  is empty also send back to network
	// buf := new(bytes.Buffer)
	// _, err := io.Copy(buf, r) -> creating the issue that buf is empty
	// if err != nil {
	// 	return err
	// }

	// tee := io.TeeReader(r, buf)

	p := &DataMessage{
		Key:  key,
		Data: buf.Bytes(),
	}

	fmt.Println(buf.Bytes())

	return s.broadcast(&Message{
		From: s.Tr
	})
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p

	log.Printf("connected with remote %s\n", p.RemoteAddr())
	return nil
}

func (s *FileServer) loop() {

	defer func() {
		log.Println("file server stopped due to user quit action")
		s.Transport.Close()
	}()

	for {
		select {
		case msg := <-s.Transport.Consume():
			// fmt.Println(msg)
			var p DataMessage
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&p); err != nil {
				log.Fatal(err)
			}
			fmt.Println("recv msg")
			fmt.Printf("%+v\n", string(p.Data))

		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) handleMessage(p *DataMessage) error {

}

func (s *FileServer) bookstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		// checking for empty
		if len(addr) == 0 {
			continue
		}

		fmt.Println("attemping to connect the remote network: ", addr)
		// s.Transport.Dial()
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error: ", err)

				// continue
			}
		}(addr)
	}

	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	if len(s.BootstrapNodes) != 0 {

		s.bookstrapNetwork()
	}
	s.loop()

	return nil
}

// func (s *FileServer) Store(key string, r io.Reader) error {
// 	return s.store.Write(key, r)
// }
