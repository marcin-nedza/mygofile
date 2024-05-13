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
	return t.rpcch
}

func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
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
	for {
		rpc := RPC{}
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			return
		}

		rpc.From=conn.RemoteAddr().String()

		t.rpcch <- rpc
	}
}
