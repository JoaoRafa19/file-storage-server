package p2p


// HandshakeFunc is like as websocket upgrader?
// 
type HandShakeFunc func(Peer) error


func NOPHandshakeFunc(Peer) error {
	return nil
}