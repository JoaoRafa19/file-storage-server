package p2p

import "net"

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

func WithOnPeer(onpeer func(Peer) error) TCPOpts {
	return func(t *TCPTransport) {
		t.onPeer = onpeer
	}
}
