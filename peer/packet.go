package peer

import (
	"bytes"
	"encoding/binary"
	"io"
	"reflect"
)

//Packet represents message received/send from/to the network
type Packet interface {
	Decode(reader io.Reader) error
	Encode() ([]byte, error)
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

	// keepAlive is message with length 0 and no id
	keepAlaiveMassage := []byte{0, 0, 0, 0}
	if reflect.DeepEqual(lenBuf, keepAlaiveMassage) {
		p.id = keepAlaive
		return nil
	}
	idBuf := make([]byte, 1)

	_, err = io.ReadFull(reader, idBuf)
	if err != nil {
		return err
	}
	p.id = idBuf[0]

	// Payload length is the total length withoud the "id" byte
	numBytes := binary.BigEndian.Uint32(lenBuf) - 1
	payload := make([]byte, numBytes)
	_, err = io.ReadFull(reader, payload)
	if err != nil {
		return err
	}
	p.payload = payload
	return nil
}

func (p *packet) Encode() ([]byte, error) {
	var buffer bytes.Buffer
	lenBuf := make([]byte, 4)
	payloadLen := len(p.Payload())
	// Total length is payload length + 1 (for id byte)
	binary.BigEndian.PutUint32(lenBuf, uint32(1+payloadLen))
	_, err := buffer.Write(lenBuf)
	if err != nil {
		return nil, err
	}
	_, err = buffer.Write([]byte{p.id})
	if err != nil {
		return nil, err
	}
	if payloadLen != 0 {
		buffer.Write(p.Payload())
	}
	return buffer.Bytes(), nil
}

func (p *packet) ID() byte {
	return p.id
}
