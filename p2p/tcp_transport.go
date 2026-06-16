package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {
	// the underlying connection of the peer
	conn net.Conn
	//if dial and accept a connection => outbound == true
	// if accpet and retrieve a connection => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// close implements the Peer interface
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransport struct {
	listenAddr string
	listener   net.Listener
	shakeHands HandshakeFunc
	decoder    Decoder
	rpcch      chan RPC
	onPeer     func(Peer) error

	mu sync.RWMutex // Mutex above the thing you want to protect
}

func NewTCPTransport(opts ...TCPOpts) *TCPTransport { //Return the struct instead of the interface for testing purposes
	transport := &TCPTransport{
		rpcch: make(chan RPC),
	}

	for _, opt := range opts {
		opt(transport)
	}

	return transport
}

// Consume implements the Transport Interface, which will return read-only channel
// for reading incomming messages received form another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (tr *TCPTransport) ListenAndAccept() error {
	var err error
	tr.listener, err = net.Listen("tcp", tr.listenAddr)
	if err != nil {
		return err
	}

	go tr.strartAcceptLoop()

	return err
}

func (tr *TCPTransport) strartAcceptLoop() {
	for {
		conn, err := tr.listener.Accept()
		if err != nil {
			fmt.Println("TCP accept error: ", err.Error())
		}

		fmt.Printf("new incomming connection %+v\n", conn)

		go tr.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error
	defer func() {
		fmt.Printf("Droping peer connection: %s\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)
	if err := t.shakeHands(peer); err != nil {
		return
	}

	if t.onPeer != nil {
		if err = t.onPeer(peer); err != nil {
			return
		}
	}

	rpc := RPC{}
	//Read Loop
	for {
		err := t.decoder.Decode(conn, &rpc)
		if err == net.ErrClosed {
			return
		}

		if err != nil {
			fmt.Printf("TCP read error: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc

		fmt.Printf("%+v\n", rpc)
	}

}
