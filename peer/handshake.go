package peer

import (
	"bytes"
	"encoding/binary"
	"errors"

	i "github.com/bkolad/gTorrent/init"
	"github.com/bkolad/gTorrent/torrent"
)

const protocol = "BitTorrent protocol"

//Bittorent protocol handshake
type Handshake struct {
	len      uint8
	protocol string
	rsvd     []byte
	InfoHash []byte
	peerID   string
}

func NewHandshake(conf i.Configuration, info *torrent.Info) Handshake {
	var h Handshake
	h.len = uint8(len(protocol))
	h.protocol = protocol
	h.rsvd = make([]byte, 8)
	h.InfoHash = info.InfoHash
	h.peerID = conf.PeerID
	return h
}

//Encode serializes Hnadshake to byte array
func (h *Handshake) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.BigEndian, h.len); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, []byte(h.protocol)); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, h.rsvd); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, h.InfoHash); err != nil {
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, []byte(h.peerID)); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

//Decode converts byte array to Handshake struct
func (h *Handshake) Decode(data []byte) error {
	if len(data) != 68 {
		return errors.New("Handshake must be 68 bytes long")
	}
	h.len = uint8(data[0])
	h.protocol = string(data[1:20])
	if h.protocol != protocol {
		return errors.New("Only BitTorrent protocol is supported")
	}
	h.rsvd = data[20:28]
	h.InfoHash = data[28:48]
	h.peerID = string(data[48:68])
	return nil
}
