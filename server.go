package main

import (
	"fmt"
	"log"
	"sync"

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

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s* FileServer) OnPeer(peer p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	
	s.peers[peer.RemoteAddr().String()] = peer
	
	log.Printf("New peer connected: %s", peer.RemoteAddr() )

	return nil
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due to quit action")
		s.Transport.Close()
	}()
	for {

		select {
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)
		case <-s.quitch:
			return
		}
	}
}
