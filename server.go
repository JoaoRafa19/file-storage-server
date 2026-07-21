package main

import (
	"bytes"
	"diststorage/p2p"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Trasport          p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	store  *Store
	quitCh chan struct{}
}

type Message struct {
	From    string
	Payload any
}

type DataMessage struct {
	Key  string
	Data []byte
}

func (s *FileServer) broadcast(p *Message) error {

	peers := []io.Writer{}

	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)

	return gob.NewEncoder(mw).Encode(p)
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	// 1 - store data in disk
	// 2 - broadcast data to all known peer in the network
	// 3 -
	//
	//
	buf := new(bytes.Buffer)

	tee := io.TeeReader(r, buf)
	//
	if err := s.store.Write(key, tee); err != nil {
		return err
	}

	p := &DataMessage{
		Key:  key,
		Data: buf.Bytes(),
	}

	return s.broadcast(&Message{
		From:    s.Trasport.ListenAddr(),
		Payload: p,
	})
}

// Starts the file serving listener and
// reads the connection on a loop returning a error
func (s *FileServer) Start() error {
	if err := s.Trasport.ListenAndAccept(); err != nil {
		log.Printf("Start error :%+v", err)
		return err
	}

	s.bootstrapNetwork()
	s.loop()

	return nil
}

func (s *FileServer) Stop() {
	close(s.quitCh)

}

func (s *FileServer) handleMessage(m *Message) error {
	switch v := m.Payload.(type) {
	case *DataMessage:
		fmt.Printf("received data %+v\n", v)
	}

	return nil
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped due user quit action")
		s.Trasport.Close()
	}()

	for {
		select {
		case msg := <-s.Trasport.Consume():
			fmt.Println("receive message")
			var m Message
			if err := gob.NewDecoder(bytes.NewReader(msg.Payload)).Decode(&m); err != nil {
				log.Fatal(err)
			}

			if err := s.handleMessage(&m); err != nil {
				log.Println(err)
			}

		case <-s.quitCh:
			return
		}
	}
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[p.RemoteAddr().String()] = p

	log.Printf("connected with remote %s", p.RemoteAddr())

	return nil
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {

		if len(addr) == 0 {
			continue
		}

		fmt.Println("Attempting to connect with remote: ", addr)
		// Server dial the peer
		go func(addr string) {
			if err := s.Trasport.Dial(addr); err != nil {
				log.Println("dial error ", err)
			}
		}(addr)
	}

	return nil
}

// NewFileServer creates a [FileServer] configured with the given options.
//
// Available options:
//   - StorageRoot the storage root path folder
//   - [PathTransformFunc] the func to transform the path string into a CAS
//   - Transport the transport to broadcast the files either be TCP, HTTP, FTP, SFTP, FSTPS, GRPC, Websocket or anything like that
//
// It returns the concrete *[FileServer] (rather than an interface) to make
// testing easier.
func NewFileServer(opts FileServerOpts) *FileServer {

	return &FileServer{
		FileServerOpts: opts,
		quitCh:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
		store: NewStore(
			WithPathTransformFunc(opts.PathTransformFunc),
			WithRootPath(opts.StorageRoot),
		),
	}
}
