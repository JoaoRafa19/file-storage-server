package main

import (
	"log"

	"github.com/JoaoRafa19/file-storage-server/p2p"
)

func main() {
	tr := p2p.NewTCPTransport(":9000")
	log.Println("Running at 9000")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
