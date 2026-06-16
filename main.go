package main

import (
	"diststorage/p2p"
	"fmt"
	"log"
)

func OnPeer(peer p2p.Peer) error {
	fmt.Println("doing some logic with the peer outside the transport")
	return nil
}

func main() {
	tr := p2p.NewTCPTransport(
		p2p.WithListenAddr(":3000"),
		p2p.WithShakeHands(p2p.NOPHandshakeFunc),
		p2p.WithDecoder(p2p.DefaultDecoder{}),
		p2p.WithOnPeer(OnPeer),
	)

	// go func() {
	// 	for {
	// 		msg := <-tr.Consume()
	// 		fmt.Printf("%+v\n", msg)
	// 	}
	// }()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
