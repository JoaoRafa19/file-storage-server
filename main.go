package main

import (
	"diststorage/p2p"
	"log"
)

func main() {
	tr := p2p.NewTCPTransport(
		p2p.WithListenAddr(":3000"),
		p2p.WithShakeHands(p2p.NOPHandshakeFunc),
		p2p.WithDecoder(p2p.DefaultDecoder{}),
	)

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
