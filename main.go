package main

import (
	"log"

	"github.com/JoaoRafa19/file-storage-server/p2p"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListenAddrs:   listenAddr,
		HandShakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		Transport:         tcpTransport,
		PathTransformFunc: CASPathTransformFunc,
		BootstapNodes:     nodes,
	}

	s := NewFileServer(fileServerOpts)
	tcpTransport.OnPeer = s.OnPeer

	return s

}

func main() {
	s1 := makeServer(":3030", "")
	s2 := makeServer(":", ":3030")

	go func() {
		log.Fatal(s1.Start())
	}()

	s2.Start()

}
