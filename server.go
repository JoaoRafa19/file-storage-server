package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/JoaoRafa19/file-storage-server/p2p"
)

type FileServerOpts struct {
	ListenAddr        string
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	TransportOpts     p2p.TCPTransportOpts
	BootstapNodes     []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer
	store    *Store
	quitch   chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := &StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(*storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

func (s *FileServer) broadCast(msg *Message) error {
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)

	return gob.NewEncoder(mw).Encode(msg)
}

type Message struct {
	From    string
	Payload any
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	//1. Store the file to disk

	// 2. broadcast the file to all know peers in the network
	// buf := new(bytes.Buffer)
	// tee := io.TeeReader(r, buf)

	// if err := s.store.Write(key, tee); err != nil {
	// 	return err
	// }

	// p := &DataMessage{
	// 	Key:  key,
	// 	Data: buf.Bytes(),
	// }

	// return s.broadCast(&Message{
	// 	From:    "todo",
	// 	Payload: p,
	// })

	buff := new(bytes.Buffer)

	msg := Message{
		Payload: []byte("storagekey"),
	}

	if err := gob.NewEncoder(buff).Encode(msg); err != nil {
		return err
	}

	for _, peer := range s.peers {
		if err := peer.Send(buff.Bytes()); err != nil {
			return err
		}
	}

	time.Sleep(time.Second * 3)

	payload := []byte("THIS LARGE FILE")
	for _, peer := range s.peers {
		if err := peer.Send(payload); err != nil {
			return err
		}
	}

	return nil

}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	if len(s.BootstapNodes) != 0 {
		s.bootstrapNetwork()
	}
	s.loop()

	return nil
}

func (s *FileServer) OnPeer(peer p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[peer.RemoteAddr().String()] = peer

	log.Printf("New peer connected: %s", peer.RemoteAddr())

	return nil
}



func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to quit action")
		s.Transport.Close()
	}()
	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Println(err)
			}

			peer, ok := s.peers[rpc.From]
			
			if !ok {
				panic("peer not found in peers")
			}

			b := make([]byte, 1000)
			if _, err := peer.Read(b); err != nil {
				panic(err)
			} 

			panic("")  

			fmt.Printf("Received: %s\nTeste", string(msg.Payload.([]byte)))
			// if err := s.handleMessage(&m); err != nil {
			// 	 log.Println(err)
			// }


		case <-s.quitch:
			return
		}
	}
}

// func (s *FileServer) handleMessage(msg *Message) error {
	
// 	return nil
// }

func (s *FileServer) bootstrapNetwork() error {

	for _, addr := range s.BootstapNodes {
		if len(addr) == 0 {
			continue
		}
		log.Println("attempting to conect the remote: ", addr)
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {

				log.Println("dial error: ", err)
			}
			fmt.Println("Conect to addr: ", addr)
		}(addr)

	}
	return nil
}
