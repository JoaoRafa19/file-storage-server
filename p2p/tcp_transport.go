package p2p

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {
	// the underlying connection of the peer. Which in this case is a TCP connection
	net.Conn
	//if dial and accept a connection => outbound == true
	// if accpet and retrieve a connection => outbound == false
	outbound bool

	Wg *sync.WaitGroup
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.Conn.Write(b)
	return err
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
		Wg:       &sync.WaitGroup{},
	}
}

type TCPTransport struct {
	listenAddr string
	listener   net.Listener
	shakeHands HandshakeFunc
	decoder    Decoder
	rpcch      chan RPC
	OnPeer     func(Peer) error

	mu sync.RWMutex // Mutex above the thing you want to protect
}

// Dial implements [Transport].
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	// handle conn
	t.handleConn(conn, true)

	return nil
}

// Consume implements the Transport Interface, which will return read-only channel
// for reading incomming messages received form another peer in the network
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (tr *TCPTransport) ListenAndAccept() error {
	var err error
	fmt.Println("Starting tcp transport")
	tr.listener, err = net.Listen("tcp", tr.listenAddr)
	if err != nil {
		fmt.Println("error on start tcp connection")
		return err
	}

	go tr.strartAcceptLoop()

	log.Printf("TCP transport listen in port %s/n", tr.listenAddr)

	return err
}

// Close implements [Transport] interface
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

func (tr *TCPTransport) strartAcceptLoop() {
	for {
		conn, err := tr.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Println("TCP accept error: ", err.Error())
		}

		go tr.handleConn(conn, false)
	}
}

func (t *TCPTransport) ListenAddr() string {
	return t.listenAddr
}

func (t *TCPTransport) handleConn(conn net.Conn, outbound bool) {
	var err error
	defer func() {
		fmt.Printf("Droping peer connection: %s\n", err)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, true)
	if err := t.shakeHands(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	rpc := RPC{}
	//Read Loop
	for {
		err := t.decoder.Decode(conn, &rpc)
		if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
			return
		}

		if err != nil {
			fmt.Printf("TCP read error: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr().String()
		peer.Wg.Add(1)
		fmt.Println("waiting stream")
		t.rpcch <- rpc
		peer.Wg.Wait()

		fmt.Println("stream continue")

	}

}

// NewTCPTransport creates a TCPTransport configured with the given options.
//
// Available options:
//   - WithListenAddr(string): address the transport listens on
//   - WithListener(net.Listener): use a pre-existing listener
//   - WithShakeHands(HandshakeFunc): handshake performed on new connections
//   - WithDecoder(Decoder): decoder used to read incoming RPC messages
//   - WithOnPeer(func(Peer) error): callback invoked when a peer connects
//
// It returns the concrete *TCPTransport (rather than an interface) to make
// testing easier.
func NewTCPTransport(opts ...TCPOpts) Transport {
	transport := &TCPTransport{
		rpcch: make(chan RPC),
	}

	for _, opt := range opts {
		opt(transport)
	}

	return transport
}
