package main

import (
	"bytes"
	"myfile/p2p"
	"time"
)

func makeServer(addr string, nodes ...string) *FileServer {
	tcpTransportOpt := p2p.TCPTransportOpts{
		ListenAddr: addr,
		Decoder:    p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpt)

	fileserverOpts := FileServerOpts{
		Transport:         tcpTransport,
		PathtransformFunc: CASPathtransformFunc,
		StorageRoot:       addr + "_storage",
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileserverOpts)
	tcpTransport.OnPeer = s.OnPeer
	return s
}
func main() {
	s1 := makeServer(":3000")
	s3 := makeServer(":3001")
	s2 := makeServer(":4000", ":3000",":3001")
	time.Sleep(time.Second)
	go func() {
		s1.Start()
	}()
	go func() {
		s3.Start()
	}()
	time.Sleep(time.Second)
	go s2.Start()
	data := bytes.NewReader([]byte("jojojojo"))
	s2.Store("heja", data)
	time.Sleep(time.Second)
	// RenderMenu()
	select {}
}
