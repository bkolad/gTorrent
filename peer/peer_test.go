package peer

import (
	"testing"

	p "github.com/bkolad/gTorrent/piece"
	"github.com/bkolad/gTorrent/torrent"
	"github.com/stretchr/testify/require"
)

func torrents() []torrent.Info {
	chunkSize := 32
	return []torrent.Info{
		torrent.Info{
			PieceSize: 15 * chunkSize,
			Length:    162 * chunkSize,
			ChunkSize: chunkSize,
		},
		torrent.Info{
			PieceSize: 15 * chunkSize,
			Length:    165 * chunkSize,
			ChunkSize: chunkSize,
		},
	}
}

func TestPeer(t *testing.T) {
	peerInfo := torrent.PeerInfo{IP: "SOME IP", Port: 9912}
	handshake := Handshake{}

	for _, torrentInfo := range torrents() {
		pieceManager := p.NewManager(torrentInfo)

		repo := makeRepo(torrentInfo.Length, torrentInfo.PieceSize)
		fakeNet := fakeNetwork(repo)
		peer := newPeerWithNetwork(fakeNet, make(chan MSG), peerInfo, handshake, pieceManager)
		pieces := make([]bool, 16)
		pieces[3] = true
		pieces[4] = true
		pieces[9] = true
		pieces[10] = true

		bitfield := bitsToBytes(pieces)

		peer.onBitfield(bitfield)
		peer.onUnchoke()

		timeout := 1000
		done := false
		for !done {
			req, payload := fakeNet.payload()
			done = peer.onPiece(req.piece, req.offset, payload)
			timeout--
			if timeout <= 0 {
				require.Fail(t, "Test Timeout")
			}
		}
		require.Equal(t, peer.pieceRepository.Get(9, 10, 10), repo.Get(9, 10, 10))
	}
}

type req struct {
	piece  uint32
	offset uint32
	size   uint32
}

type fakeNet struct {
	repo      p.Repository
	requested req
}

func (fN *fakeNet) payload() (req, []byte) {
	req := fN.requested
	p := fN.repo.Get(req.piece, req.offset, req.size)
	return req, p
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

func fakeNetwork(repo p.Repository) *fakeNet {
	return &fakeNet{repo, req{}}
}

func makeRepo(length, pieceSize int) p.Repository {
	lastPieceSize, numberOfPieces := p.CalculateLastPieceSize(length, pieceSize)
	repo := p.NewRepo(numberOfPieces)
	for i := uint32(0); i < numberOfPieces-1; i++ {
		data := make([]byte, pieceSize)
		for k := 0; k < pieceSize; k++ {
			data[k] = byte(uint32(3*k) + 2*i)
		}
		repo.Save(uint32(i), data)
	}

	piece := make([]byte, lastPieceSize)
	for k := uint32(0); k < lastPieceSize; k++ {
		piece[k] = byte(k)
	}

	repo.Save(uint32(numberOfPieces-1), piece)
	return repo
}
