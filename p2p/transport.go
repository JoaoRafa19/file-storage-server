package p2p

// Peer is a interface that represents the remote node
type Peer interface {
}

// Transport is anything that handles the comunication
// between the nodes in the network. This can be of the form
// tcp, udp, websockets ...
type Transport interface {
	ListenAndAccept() error
}
