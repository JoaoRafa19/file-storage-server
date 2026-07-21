package p2p

import "net"

const StreamType = 0x2
const MessageType = 0x1

// RPC represents any arbitrary data that is being send
// over each transport between two nodes
type RPC struct {
	Payload []byte
	From    net.Addr
}
