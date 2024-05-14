package p2p

import (
	"fmt"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type DefaultDecoder struct {
}

type DefaultDecoder2 struct {
}

func (dec DefaultDecoder) Decode(r io.Reader, msg *RPC) error {
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	msg.Payload = buf[:n]

	 // fmt.Printf("--------Decode:%+v ",string(msg.Payload))
	return nil
}

func (dec DefaultDecoder2) Decode(r io.Reader, msg *RPC) error {
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	msg.Payload = buf[:n]
	fmt.Println("FANCY msg", string(msg.Payload))

	return nil
}
