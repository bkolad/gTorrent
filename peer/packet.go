package peer

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"

	log "github.com/bkolad/gTorrent/logger"
)

type Packet interface {
	Decode(reader io.Reader) error
	Encode() []byte
	ID() byte
	Payload() []byte
}

type packet struct {
	id      byte
	payload []byte
}

func (p *packet) Payload() []byte {
	return p.payload
}

func (p *packet) Decode(reader io.Reader) error {
	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(reader, lenBuf)
	if err != nil {
		//TODO handel EOF (remote peer closes connection)
		return err
	}

	keepAlaiveMassage := []byte{0, 0, 0, 0}
	if reflect.DeepEqual(lenBuf, keepAlaiveMassage) {
		log.Info("Keep alive recieved")
		return nil
	}
	idBuf := make([]byte, 1)

	_, err = io.ReadFull(reader, idBuf)
	if err != nil {
		return err
	}

	numBytes := binary.BigEndian.Uint32(lenBuf)
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
	var b bytes.Buffer
	lenBuf := make([]byte, 4)
	payloadLen := len(p.Payload())
	binary.BigEndian.PutUint32(lenBuf, uint32(1+payloadLen))
	b.Write(lenBuf)
	b.Write([]byte{p.id})
	if payloadLen != 0 {
		b.Write(p.Payload())
	}

	return b.Bytes()
}

func (p *packet) ID() byte {
	return p.id
}
