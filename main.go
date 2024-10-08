package main

import (
	"bytes"
	"log"
	"time"

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
		StorageRoot:       "storage",
		Transport:         tcpTransport,
		PathTransformFunc: CASPathTransformFunc,
		BootstapNodes:     nodes,
	}

	s := NewFileServer(fileServerOpts)
	tcpTransport.OnPeer = s.OnPeer

	return s

}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.Start())
	}()

	time.Sleep(time.Second * 4)
	go s2.Start()
	time.Sleep(time.Second * 4)

	data := bytes.NewReader([]byte("big data here"))
	s2.StoreData("myprivatedata", data)

	select {}
}
