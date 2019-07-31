package peer

import (
	"bytes"
	"encoding/binary"
)

//type msgId int

const (
	keepAlaive = iota - 1
	choke
	unchoke
	interested
	notInterested
	have
	bitfield
	request
	piece
	cancel
	port
	unknown
)

func decodePiece(payload []byte) (uint32, uint32, []byte) {
	piece := binary.BigEndian.Uint32(payload[0:4])
	offset := binary.BigEndian.Uint32(payload[4:8])
	return piece, offset, payload[8:]
}

func encodePieceRequest(piece, offset, size uint32) Packet {
	b := &bytes.Buffer{}
	binary.Write(b, binary.BigEndian, piece)
	binary.Write(b, binary.BigEndian, offset)
	binary.Write(b, binary.BigEndian, size)
	packet := &packet{id: request, payload: b.Bytes()}
	return packet
}

func encodeHave(bitfield []byte) Packet {
	return &packet{id: have, payload: bitfield}
}

func encodeInterested() Packet {
	return &packet{id: interested, payload: nil}
}
