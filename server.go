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
		store:          NewStore(storeOpts),
	}
}

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
	if err := s.store.writeStream(key, teebuf); err != nil {
		return err
	}
	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: 10,
		},
	}
	time.Sleep(time.Second)
	//broadcast to all other connected peers
	if err := s.broadcast(&msg); err != nil {
		return err
	}
	time.Sleep(time.Second * 2)
	for _, peer := range s.peers {
		n, err := io.Copy(peer, filebuf)
		if err != nil {
			return err
		}
		fmt.Printf("received and written (%d) bytes", n)
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
		if len(addr) == 0 {
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
			fmt.Printf("===rpc: %+v", string(rpc.Payload))
			if err := gob.NewDecoder(bytes.NewBuffer(rpc.Payload)).Decode(&msg); err != nil {
				fmt.Println("Error decoding:", err)
				return
			}

			fmt.Printf("\n===rpc---msg: %+v", msg)
			if err := s.handleMessage(rpc.From, &msg); err != nil {
				log.Println("Error handling message: ", err)
			}
		}

	}
}

func (s *FileServer) handleStoreFile(from string, msg MessageStoreFile) error {
	var (
		buf = new(bytes.Buffer)
	)
	peer, ok := s.peers[from]

	if !ok {
		return fmt.Errorf("peer (%s) could not be found in the peer list", from)
	}

	_, err := io.Copy(buf, peer)
	if err != nil {
		return err
	}
	fmt.Printf("\n\t\t==== %+v", msg)
	if err := s.store.writeStream(msg.Key, buf); err != nil {
		fmt.Println("Error: ", err)
	}

	// peer.(*p2p.TCPPeer).Wg.Done()
	return nil
}

func (s *FileServer) handleMessage(from string, msg *Message) error {
	fmt.Println("---handling message")
	switch v := msg.Payload.(type) {
	case MessageStoreFile:

		fmt.Printf("---handling message:---%+v", v)
		return s.handleStoreFile(from, v)
	}

	return nil
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
	gob.Register(MessageStoreFile{})
}
