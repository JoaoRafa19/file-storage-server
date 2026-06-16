package p2p

import "net"

// RPC represents any arbitrary data that is being send
// over each transport between two nodes
type RPC struct {
	Payload []byte
	From    net.Addr
}
