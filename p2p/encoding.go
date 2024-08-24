package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *Message) error
}

type GOBDecoder struct{}

func (g GOBDecoder) Decode(r io.Reader, v*Message) error {
	return gob.NewDecoder(r).Decode(v)
}

type DefaultDecoder struct{}

func (g DefaultDecoder) Decode(r io.Reader, v *Message) error {
	buff := make([]byte, 1028)
	n, err := r.Read(buff)
	if err != nil {
		return err
	}
	v.Payload = buff[:n]

	return nil
}
