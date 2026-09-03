package p2p

import (
	// "bytes"
	"encoding/gob"
	// "fmt"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *Message) error
}

type GOBDecoder struct{}

func (dec GOBDecoder) Decode(r io.Reader, msg *Message) error {
	return gob.NewDecoder(r).Decode(msg)
}

type DefaultDecoder struct{}

func (dec DefaultDecoder) Decode(r io.Reader, msg *Message) error {
	// buf := new(bytes.Buffer)
	// buf.Read(r)
	buf := make([]byte, 1028)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	// fmt.Println(string(buf[:n]))
	msg.Payload = buf[:n]

	return nil
}
