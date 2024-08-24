package main

import (
	"log"

	"github.com/JoaoRafa19/file-storage-server/p2p"
)

func main() {
	tcpOpts := p2p.TCPTransportOpts{ 
		ListenAddrs: ":9000",
		Decoder:  p2p.DefaultDecoder{},
		HandShakeFunc:  p2p.NOPHandshakeFunc,
	}
	tr := p2p.NewTCPTransport(tcpOpts)
	log.Println("Running at 9000 !!!!")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
