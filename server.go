package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"myfile/p2p"
	"sync"
	"time"
)

type FileServerOpts struct {
	Transport         p2p.Transport
	PathtransformFunc PathtransformFunc
	StorageRoot       string
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts
	onPeerLock sync.Mutex
	peers      map[string]p2p.Peer

	store *Store
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathtransformFunc: opts.PathtransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		peers:          make(map[string]p2p.Peer),
		store:          NewStore(storeOpts),
	}
}

func (s *FileServer) broadcast(msg *Message) error {
	buf := new(bytes.Buffer)
	//encode data
	fmt.Printf("broadcasting message to all peers\n")
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}
	for _, peer := range s.peers {
		fmt.Println("Sending from loop", peer.LocalAddr())
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}
	//send data to all peers
	return nil
}

func (s *FileServer) Store(key string, r io.Reader) error {
	var (
		filebuf = new(bytes.Buffer)
		teebuf  = io.TeeReader(r, filebuf)
	)
	fmt.Println("Attempting to store key")
	if err := s.store.writeStream(key, teebuf); err != nil {
		return err
	}
	msg := Message{
		Payload: key,
	}
	time.Sleep(time.Second)
	//broadcast to all other connected peers
	if err := s.broadcast(&msg); err != nil {
		return err
	}
	return nil
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.onPeerLock.Lock()
	defer s.onPeerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p
	log.Printf("connected with remote %s\n", p.RemoteAddr())
	return nil
}

type Message struct {
	Payload any
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr)==0{
			continue
		}
		fmt.Printf("Attempting to connect with remote [%s]\n", addr)
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				log.Println(err)
			}
		}(addr)
	}
	return nil
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
			fmt.Printf("Decoded msg %+v\n", msg)
		}

	}
}

func (s *FileServer) Start() error {
	fmt.Println("Server start")
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	//bootstrap nodes
	s.bootstrapNetwork()
	//start server loop
	s.loop()

	return nil
}
func init() {
	gob.Register(Message{})
}
