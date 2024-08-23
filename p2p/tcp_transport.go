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

type TCPTransport struct {
	listenAddress string
	listener      net.Listener
	shakeHands    HandShakeFunc // like an websocket upgrader
	decoder       Decoder

	mu    sync.Mutex
	peers map[net.Addr]Peer
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		shakeHands:    NOPHandshakeFunc,
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.listenAddress)
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

type Temp struct {}

func (t *TCPTransport) handleConn(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.shakeHands(peer); err != nil {
	}
	
	
	// buff := new(bytes.Buffer)
	// 
	msg:= &Temp{}
	for {
		// n,err := conn.Read(buff)
		//
		if err := t.decoder.Decode(conn, msg); err !=nil {
			fmt.Printf("TCP ERROR: %s\n", err)
			continue
		}
	}

}
