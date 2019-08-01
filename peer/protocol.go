package peer

import (
	"bytes"
	"encoding/binary"
)

const (
	choke = iota
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
	// keepAlaive is a special message in the protocol and it doesn't have message id
	// here we reserve 255 for it, and handle it like any other message.
	keepAlaive = 255
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

func decodeRequest(payload []byte) (uint32, uint32, uint32) {
	piece := binary.BigEndian.Uint32(payload[0:4])
	offset := binary.BigEndian.Uint32(payload[4:8])
	size := binary.BigEndian.Uint32(payload[8:12])
	return piece, offset, size
}

func encodeHave(bitfield []byte) Packet {
	return &packet{id: have, payload: bitfield}
}

func encodeInterested() Packet {
	return &packet{id: interested, payload: nil}
}
