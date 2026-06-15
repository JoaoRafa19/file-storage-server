package p2p

import "net"

// Message represents any arbitrary data that is being send
// over each transport between two nodes
type Message struct {
	Payload []byte
	From    net.Addr
}
