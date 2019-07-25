package peer

import (
	"encoding/binary"
	"io"
)

type Packet interface {
	Decode(reader io.Reader) error
	Encode() []byte
	ID() byte
}

type packet struct {
	id      byte
	payload []byte
}

func (p *packet) Decode(reader io.Reader) error {
	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(reader, lenBuf)
	if err != nil {
		//TODO handel EOF (remote peer closes connection)
		return err
	}
	numBytes := binary.BigEndian.Uint32(lenBuf)
	idBuf := make([]byte, 1)

	_, err = io.ReadFull(reader, idBuf)
	if err != nil {
		return err
	}

	size := numBytes - 1
	payload := make([]byte, size)
	_, err = io.ReadFull(reader, payload)
	if err != nil {
		return err
	}
	p.id = idBuf[0]
	p.payload = payload
	return nil
}

func (p *packet) Encode() []byte {
	return nil
}

func (p *packet) ID() byte {
	return p.id
}
