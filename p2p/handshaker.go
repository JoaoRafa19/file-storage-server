package p2p

// Handshakefunc is
type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error { return nil }
