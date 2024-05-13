package main

import (
	"myfile/p2p"
	"time"
)

func makeServer(addr string) *FileServer{
	tcpTransportOpt := p2p.TCPTransportOpts{
		ListenAddr: addr,
		Decoder: p2p.DefaultDecoder2{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpt)

	fileserverOpts := FileServerOpts{
		Transport:  tcpTransport,
	}

	s :=  NewFileServer(fileserverOpts)
	tcpTransport.OnPeer=s.OnPeer	
	return s
}
func main() {
	s1 := makeServer(":3000")
	time.Sleep(time.Second)
	go s1.Start()
	time.Sleep(time.Second)
	select {}
}
