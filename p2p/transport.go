package p2p

// Peer is a interface thar representes the remote node.
type Peer interface {
    
}

// Transport is anything that handles comunication
// between nodes in network. This can be in the form of (TCP, websockets, ...)
type Transport interface {
	ListenAndAccept() error    
}