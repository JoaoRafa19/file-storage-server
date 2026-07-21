package main

import (
	"bytes"
	"diststorage/p2p"
	"fmt"
	"log"
	"time"
)

func OnPeer(peer p2p.Peer) error {
	fmt.Println("doing some logic with the peer outside the transport")
	return nil
}

func makeServer(listenAddr string, nodes ...string) *FileServer {

	fileserverOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		BootstrapNodes:    nodes,
	}

	s := NewFileServer(fileserverOpts)
	tcpTransport := p2p.NewTCPTransport(
		p2p.WithListenAddr(listenAddr),
		p2p.WithShakeHands(p2p.NOPHandshakeFunc),
		p2p.WithDecoder(p2p.DefaultDecoder{}),
		p2p.WithOnPeer(s.OnPeer),
	)

	s.Trasport = tcpTransport

	return s
}

func main() {

	s1 := makeServer(":3000")
	s2 := makeServer(":4000", ":3000")

	go func() {
		if err := s1.Start(); err != nil {
			log.Fatal(err)
		}
	}()

	time.Sleep(time.Second * 3)

	go s2.Start()
	time.Sleep(time.Second * 3)

	data := bytes.NewReader([]byte("My big data file here !"))
	s2.StoreData("private_data", data)

	select {}
}
