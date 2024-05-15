package main

import (
	"log"
	"myfile/p2p"
	"time"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathtransformFunc: CASPathtransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileServerOpts)

	tcpTransport.OnPeer = s.OnPeer
	return s

}

func main() {
	s1 := makeServer(":3000", ":4000")
	s2 := makeServer(":4000", ":3000")
	// s3 := makeServer(":4001")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(time.Second * 2)
	go s2.Start()
	time.Sleep(time.Second)

	// data := bytes.NewReader([]byte("my big data titile here!"))
	// s2.Store("coolPicture.jpg", data)
	time.Sleep(time.Millisecond * 5)
	RenderMenu(s2)
	select {}
}
