package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"myfile/p2p"
	"sync"
)

type FileServerOpts struct {
	// ListenAddr string
	Transport p2p.Transport
}

type FileServer struct {
	FileServerOpts
	onPeerLock sync.Mutex
	peers      map[string]p2p.Peer
}

func NewFileServer(opts FileServerOpts) *FileServer {
	return &FileServer{
		FileServerOpts: opts,
		peers:          make(map[string]p2p.Peer),
	}
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.onPeerLock.Lock()
	defer s.onPeerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p
	log.Printf("connected with remote %s", p.RemoteAddr())
	return nil
}

type Message struct {
	Payload any
}

func (s *FileServer) loop() {

	defer func() {
		s.Transport.Close()
	}()
	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message

			if err := gob.NewDecoder(bytes.NewBuffer(rpc.Payload)).Decode(&msg); err != nil {
				fmt.Println("Error decoding:", err)
				return
			}
		}

	}
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	//start server loop
	s.loop()

	return nil
}
func init(){
	gob.Register(Message{})
}
