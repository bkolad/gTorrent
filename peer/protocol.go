package peer

import (
	"bytes"
	"encoding/binary"
)

const (
	choke byte = iota
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

func haveToIndex(bytes []byte) uint32 {
	index := binary.BigEndian.Uint32(bytes)
	return index
}

func bytesToBits(bytes []byte) []bool {
	bitSet := make([]bool, 0)
	for _, b := range bytes {
		for i := 0; i < 8; i++ {
			val := b & (1 << uint(i))
			hasBit := val > 0
			bitSet = append(bitSet, hasBit)
		}
	}
	return bitSet
}

func bitsToBytes(bits []bool) []byte {
	var bytes []byte
	for i := 0; i < len(bits)/8; i++ {
		var result byte
		oneByte := bits[8*i : 8*i+8]
		for k, b := range oneByte {
			if b {
				result = result + (1 << uint(k))
			}
		}
		bytes = append(bytes, result)
		result = 0
	}
	return bytes
}
