package p2p

// Message holds any arbitrary data that is beign sent  over the
// each transport b/w two nodes in the network.
// type Message struct {
// 	From    net.Addr
// 	Payload []byte
// }

type RPC struct {
	// From net.Addr
	From    string
	Payload []byte
}
