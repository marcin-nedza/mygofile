package p2p

import (
	"fmt"
	"net"
)

type TCPPeer struct {
	net.Conn
}

func NewTCPPeer(conn net.Conn) *TCPPeer {
	return &TCPPeer{
		Conn: conn,
	}
}

func (t *TCPPeer) Send(m []byte) error {
	 fmt.Println("Sending bytes through tcp conn")
	n, err := t.Conn.Write(m)
	if err != nil {
		return err
	}
	fmt.Printf("Wrote (%d) bytes to conn",n)
	return nil
}
func (t *TCPPeer) Close() error {
	return t.Close()
}

type TCPTransportOpts struct {
	ListenAddr string
	Decoder    Decoder
	OnPeer     func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC, 1024),
	}
}

func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	fmt.Println("conn", conn)
	return nil
}

func (t *TCPTransport) Consume() <-chan RPC {
	fmt.Println("consume")
	return t.rpcch
}

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	fmt.Println("Tcp Listen")
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		fmt.Println("Tcp Accept")
		if err != nil {
			fmt.Printf("TCP accept loop err: %s\n", err)
		}
		go t.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {

	defer func() {
		conn.Close()
	}()

	peer := NewTCPPeer(conn)
	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
	}
	fmt.Println("Tcp handle con")
	for {
		rpc := RPC{}
		fmt.Println("got somhinh")
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			return
		}
		fmt.Printf("Tcp Decoded %+v", string(rpc.Payload))
		rpc.From = conn.RemoteAddr().String()

		t.rpcch <- rpc
	}
}
