package p2p

import "net"

// Peer is a interface thar representes the remote node.
type Peer interface {
	Send([]byte) error
    Close() error 
	RemoteAddr() net.Addr
}

// Transport is anything that handles comunication
// between nodes in network. This can be in the form of (TCP, websockets, ...)
type Transport interface {
	ListenAndAccept() error   
	Consume() <-chan RPC
	Close() error
	Dial(string) error
}

