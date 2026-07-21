package p2p

import "net"

// Peer is a interface that represents the remote node
type Peer interface {
	net.Conn
	Send([]byte) error
}

// Transport is anything that handles the comunication
// between the nodes in the network. This can be of the form
// tcp, udp, websockets ...
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
	ListenAddr() string
	Dial(addr string) error
}
