package p2p

import (
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type DefaultDecoder struct {
}

func (dec DefaultDecoder) Decode(r io.Reader, msg *RPC) error {
	peekbuf := make([]byte, 1)
	if _, err := r.Read(peekbuf); err != nil {
		return nil
	}

	stream := peekbuf[0] == IncomingStream
	if stream {
		msg.Stream = true
		return nil
	}
	buf := make([]byte, 1024)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	msg.Payload = buf[:n]

	return nil
}
