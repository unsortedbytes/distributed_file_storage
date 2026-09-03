package p2p

// import "errors"

// ErrInvalidHandshake  is retured if the handshake between
// the local and remote node  could not be established.
// var ErrInvalidHandshake = errors.New("invalid handshake")

// type Handshaker interface{
// 	Handshake() error
// }

type HandshakeFunc func(Peer) error

// type DefaultHandshaker struct {

// }

func NOPHandshakeFunc(Peer) error {
	return nil
}
