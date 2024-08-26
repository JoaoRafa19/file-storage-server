package main

import (
	"fmt"
	"log"

	"github.com/JoaoRafa19/file-storage-server/p2p"
)

func  OnPeer (peer p2p.Peer ) error {
	
	fmt.Println("Doing som logic outside the TCP transport")
	return nil
}

func main() {
	tcpOpts := p2p.TCPTransportOpts{ 
		ListenAddrs: ":9000",
		Decoder:  p2p.DefaultDecoder{},
		HandShakeFunc:  p2p.NOPHandshakeFunc,
		OnPeer : OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)
	
	
	log.Println("Running at 9000 !!!!")
	
	go func(){
		for {
			msg := <- tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	

	select {}
}
