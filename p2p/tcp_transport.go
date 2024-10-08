 package p2p

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

// TCPPer represents the remote node over a TCP conections
type TCPPeer struct {
	// underlying connection of the peer 
	// TCP Connection
	net.Conn

	// if we dial a connection => outbound = true
	// if accept and retrieve a conn => outbound = false
	outbound bool

	Wg *sync.WaitGroup
}

type TCPTransportOpts struct {
	ListenAddrs   string
	HandShakeFunc HandShakeFunc // like an websocket upgrader
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	// mu       sync.Mutex
	rpcChan chan RPC
}

// Close implements the transport interface.
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

func (t *TCPPeer) Send(b []byte) error{
	_, err := t.Conn.Write(b)

	return err
}


func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		Wg: &sync.WaitGroup{}  ,
	}
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcChan:          make(chan RPC),
	}
}

// Consume implements the transport interface
// returning a read-only channel for reading incomming
// messages receives from an other peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcChan
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddrs)
	if err != nil {
		return err
	}

	go t.startAcceptLoop()
	log.Println("TCP Transport listen on port:", t.ListenAddrs)
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {

		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		go t.handleConn(conn, false)
	}
}

// Dial implements the Transport interface.
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}


	go t.handleConn(conn, true)

	return nil

}


func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error
	peer := NewTCPPeer(conn, outbound)

	defer func() {
		fmt.Printf("Drop peer connection: %v", err)
		conn.Close()
	}()

	if err := t.HandShakeFunc(peer); err != nil {
		return
	}
 
	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}
	


	rpc := RPC{}

	for {

		err := t.Decoder.Decode(conn, &rpc)

		if err != nil {
			t, isOpError := err.(*net.OpError)
			if isOpError {
				fmt.Println(t)
				return
			}

			fmt.Printf("TCP ERROR: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr().String()
		peer.Wg.Add(1 )
		fmt.Println("waiting till stream is done")
		t.rpcChan <- rpc
		peer.Wg.Wait()
		fmt.Println("Stream done continuing normal read loop" )
	}
}
