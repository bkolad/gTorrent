package peer

import (
	"testing"

	p "github.com/bkolad/gTorrent/piece"
	"github.com/bkolad/gTorrent/torrent"
)

func TestPeer(t *testing.T) {

	peerInfo := torrent.PeerInfo{IP: "SOME IP", Port: 9912}
	handshake := Handshake{}
	chunkSize := 10
	torrentInfo := torrent.Info{
		PieceSize: 15 * chunkSize,
		Length:    165 * chunkSize,
		ChunkSize: chunkSize,
	}
	pieceManager := p.NewManager(torrentInfo)

	data := data(torrentInfo.Length)
	fakeNet := fakeNetwork(data, torrentInfo.PieceSize)
	peer := newPeerWithNetwork(fakeNet, make(chan MSG), peerInfo, handshake, pieceManager)
	pieces := make([]bool, 16)
	pieces[3] = true
	pieces[4] = true
	pieces[9] = true
	pieces[10] = true

	bitfield := bitsToBytes(pieces)
	//fmt.Println(bitfield)
	//fmt.Println(bytesToBits(bitfield))
	peer.onBitfield(bitfield)
	peer.onUnchoke()

	i := 0
	done := false
	for !done {
		req, payload := fakeNet.payload()
		done = peer.onPiece(req.piece, req.offset, payload)
		i++
		if i > torrentInfo.Length/torrentInfo.ChunkSize {
			break
		}
	}

}

type req struct {
	piece  uint32
	offset uint32
	size   uint32
}

type fakeNet struct {
	pieces    []byte
	pieceSize int
	requested req
}

func (fN *fakeNet) payload() (req, []byte) {
	req := fN.requested
	//TODO handle the last piece request
	from := req.piece*uint32(fN.pieceSize) + req.offset
	to := from + req.size
	return req, fN.pieces[from:to]
}

func (fN *fakeNet) SendHandshake() error {
	return nil
}

func (fN *fakeNet) RegisterListener(Listener) {

}

func (fN *fakeNet) Send(p Packet) error {
	if p.ID() == request {
		piece, offset, size := decodeRequest(p.Payload())
		fN.requested = req{piece, offset, size}
	}
	return nil
}

func fakeNetwork(data []byte, pieceSize int) *fakeNet {
	return &fakeNet{data, pieceSize, req{}}
}

func data(length int) []byte {
	data := make([]byte, length)
	for i := 0; i < length; i++ {
		data[i] = byte(i)
	}
	return data
}
