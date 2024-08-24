package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
    listenAddr:=":4000"
    tr := NewTCPTransport(TCPTransportOpts{})
    
    assert.Equal(t, tr.ListenAddrs, listenAddr)
    
    // Server
    // 
    assert.Nil(t, tr.ListenAndAccept())
    
    select {}
}