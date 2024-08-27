package main

import (
	"log"

	"github.com/JoaoRafa19/file-storage-server/p2p"
)

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddrs:   ":9000",
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},

		OnPeer: func(p2p.Peer) error {
			return nil
		},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	fileServerOpts := FileServerOpts{
		StorageRoot:       "9000_network",
		Transport:         tcpTransport,
		PathTransformFunc: CASPathTransformFunc,
		BootstapNodes: []string{
			":4000",
		},
	}
	s := NewFileServer(fileServerOpts)

	// go func() {
	// 	time.Sleep(time.Second * 3)
	// 	s.Stop()
	// }()

	if err := s.Start(); err != nil {
		log.Fatal(err)
	}

}
