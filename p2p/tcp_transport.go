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

type TCPOpts func(*TCPTransport)

func WithListenAddr(addr string) TCPOpts {
	return func(t *TCPTransport) {
		t.listenAddr = addr
	}
}

func WithListener(listener net.Listener) TCPOpts {
	return func(t *TCPTransport) {
		t.listener = listener
	}
}

func WithShakeHands(shf HandshakeFunc) TCPOpts {
	return func(t *TCPTransport) {
		t.shakeHands = shf
	}
}

func WithDecoder(dec Decoder) TCPOpts {
	return func(t *TCPTransport) {
		t.decoder = dec
	}
}

type TCPTransport struct {
	listenAddr string
	listener   net.Listener
	shakeHands HandshakeFunc
	decoder    Decoder

	mu    sync.RWMutex // Mutex above the thing you want to protect
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts ...TCPOpts) *TCPTransport { //Return the struct instead of the interface for testing purposes
	transport := &TCPTransport{}

	for _, opt := range opts {
		opt(transport)
	}

	return transport
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

		fmt.Println("new incomming connection %+v", conn)

		go tr.handleConn(conn)
	}
}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.shakeHands(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	msg := &Message{}
	//Read Loop
	for {
		if err := t.decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error: %s\n", err)
			continue
		}
		msg.From = conn.RemoteAddr()

		fmt.Printf("message %+v\n", msg)
	}

}
