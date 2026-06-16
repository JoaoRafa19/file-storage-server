package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := ":4000"
	tr := NewTCPTransport(
		WithListenAddr(listenAddr),
		WithDecoder(DefaultDecoder{}),
		WithShakeHands(NOPHandshakeFunc),
	)

	assert.Equal(t, tr.listenAddr, ":4000")

	assert.Nil(t, tr.ListenAndAccept())

}
