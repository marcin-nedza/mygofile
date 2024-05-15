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
	ListenAddr        string
	Transport         p2p.Transport
	PathtransformFunc PathtransformFunc
	StorageRoot       string
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	store  *Store
	quitch chan struct{}
}

type MessageStoreFile struct {
	Key  string
	Size int64
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathtransformFunc: opts.PathtransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		peers:          make(map[string]p2p.Peer),
		quitch:         make(chan struct{}),
		store:          NewStore(storeOpts),
	}
}

// func (s *FileServer) broadcast(msg *Message) error {
// 	buf := new(bytes.Buffer)
// 	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
// 		return err
// 	}
//
// 	for _, peer := range s.peers {
// 		peer.Send([]byte{p2p.InocomingMessage})
// 		if err := peer.Send(buf.Bytes()); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func (s *FileServer) broadcast(msg *Message) error {
	buf := new(bytes.Buffer)
	//encode data
	fmt.Printf("broadcasting message to all peers\n")
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		fmt.Println("errrer", err)
		return err
	}
	for _, peer := range s.peers {
		fmt.Println("Sending from loop")
		peer.Send([]byte{p2p.IncomingMessage})
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}
	fmt.Println("Sent")
	//send data to all peers
	return nil
}

func (s *FileServer) Store(key string, r io.Reader) error {
	var (
		filebuf = new(bytes.Buffer)
		teebuf  = io.TeeReader(r, filebuf)
	)
	fmt.Println("Attempting to store key: ", key)
	size, err := s.store.Write(key, teebuf)
	if err != nil {
		return err
	}
	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: size,
		},
	}
	//broadcast to all other connected peers
	if err := s.broadcast(&msg); err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 5)
	for _, peer := range s.peers {
		peer.Send([]byte{p2p.IncomingStream})
		n, err := io.Copy(peer, filebuf)
		if err != nil {
			return err
		}
		fmt.Printf("received and written (%d) bytes", n)
	}
	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p
	log.Printf("connected with remote %s", p.RemoteAddr())
	return nil
}

type Message struct {
	Payload any
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		fmt.Println("attempting to connect with remote:", addr)
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial error:", err)
			}
		}(addr)
	}
	return nil
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to error or user quit action")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Println("decoding error:", err)
				return
			}
			if err := s.handleMessage(rpc.From, &msg); err != nil {
				log.Println("handle message error:", err)
				return
			}
		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) handleStoreFile(from string, msg MessageStoreFile) error {
	fmt.Printf("Message: %+v\n", msg)
	peer, ok := s.peers[from]

	if !ok {
		return fmt.Errorf("peer (%s) could not be found in the peer list", from)
	}

	n, err := s.store.Write(msg.Key, io.LimitReader(peer, msg.Size))
	if err != nil {
		return err
	}
	fmt.Printf("[%s] written %d bytes to disk\n", s.Transport.Addr(), n)
	peer.CloseStream()

	return nil
}

func (s *FileServer) handleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return s.handleStoreFile(from, v)
	}

	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.bootstrapNetwork()

	s.loop()

	return nil
}
func init() {
	gob.Register(MessageStoreFile{})
}
