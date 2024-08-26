package p2p

import "net"

//RPC hold the data that will be sent over
// the transports between two nodes in the network
type RPC struct {
	From net.Addr
	Payload []byte
}