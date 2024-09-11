package p2p


//RPC hold the data that will be sent over
// the transports between two nodes in the network
type RPC struct {
	From string
	Payload []byte
}