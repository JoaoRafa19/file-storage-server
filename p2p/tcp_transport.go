package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPer represents the remote node over a TCP conections
type TCPPeer struct {
	// conn is the connection of the Peer
	conn net.Conn

	// if we dial a connection => outbound = true
	// if accept and retrieve a conn => outbound = false
	outbound bool
}

type TCPTransportOpts struct {
	ListenAddrs   string
	HandShakeFunc HandShakeFunc // like an websocket upgrader
	Decoder       Decoder

}

type TCPTransport struct {
	TCPTransportOpts
	listener      net.Listener

	mu    sync.Mutex
	peers map[net.Addr]Peer
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddrs)
	if err != nil {
		return err
	}
	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {

		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		go t.handleConn(conn)
	}
}


func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.HandShakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error %s\n", err)
		return
	}

	msg := &Message{}

	for {
		
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP ERROR: %s\n", err)
			continue
		}
		
		msg.From = conn.RemoteAddr()

		fmt.Printf("MESSAGE %+v\n", msg)
	}
}
