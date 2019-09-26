package peer

import (
	"crypto/sha1"
	"testing"

	p "github.com/bkolad/gTorrent/piece"
	"github.com/bkolad/gTorrent/torrent"
	"github.com/stretchr/testify/require"
)

type torrentData struct {
	chunkSize int
	pieceSize int
	length    int
}

func torrents() []torrentData {
	return []torrentData{
		torrentData{chunkSize: 32, pieceSize: 15 * 32, length: 162 * 32},
		//torrentData{chunkSize: 32, pieceSize: 15 * 32, length: 165 * 32},
	}
}
func torrentsInfo(torrentData torrentData) (torrent.Info, p.Repository) {
	repo, pieceHashes := makeRepo(torrentData.length, torrentData.pieceSize)
	ti := torrent.Info{
		PieceSize:   torrentData.pieceSize,
		Length:      torrentData.length,
		ChunkSize:   torrentData.chunkSize,
		PieceHashes: pieceHashes,
	}
	return ti, repo
}

func TestPeer(t *testing.T) {
	peerInfo := torrent.PeerInfo{IP: "SOME IP", Port: 9912}
	handshake := Handshake{}

	for _, torrentData := range torrents() {
		torrentInfo, repo := torrentsInfo(torrentData)
		pieceManager := p.NewManager(torrentInfo, repo)

		fakeNet := fakeNetwork(repo)
		peersRepo := p.NewRepo(pieceManager.PieceCount())
		peer := newPeerWithNetwork(fakeNet, make(chan MSG), peerInfo, handshake, pieceManager, peersRepo)
		pieces := make([]bool, 16)
		pieces[3] = true
		pieces[4] = true
		pieces[9] = true
		pieces[10] = true

		bitfield := bitsToBytes(pieces)

		peer.onBitfield(bitfield)
		peer.onUnchoke()

		timeout := 20
		done := false
		for !done {
			req, payload, err := fakeNet.payload()
			require.Nil(t, err)
			packet := encodePieceData(req.piece, req.offset, payload)
			done = peer.NewPacket(packet)
			timeout--
			if timeout <= 0 {
				require.Fail(t, "Test Timeout")
			}
		}
		peerPiece, err := peer.pieceRepository.Get(9)
		require.Nil(t, err)
		piece, err := repo.Get(9)
		require.Nil(t, err)
		require.Equal(t, peerPiece, piece)
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

func (fN *fakeNet) payload() (req, []byte, error) {
	req := fN.requested
	p, err := fN.repo.Get(req.piece)
	return req, p, err
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

func makeRepo(length, pieceSize int) (p.Repository, [][]byte) {
	pieceHashes := [][]byte{}
	lastPieceSize, numberOfPieces := torrent.CalculateLastPieceSize(length, pieceSize)
	repo := p.NewRepo(uint32(numberOfPieces))
	for i := uint32(0); i < uint32(numberOfPieces-1); i++ {
		data := make([]byte, pieceSize)
		for k := 0; k < pieceSize; k++ {
			data[k] = byte(uint32(3*k) + 2*i)
		}
		pieceHash := sha1.Sum(data)
		pieceHashes = append(pieceHashes, pieceHash[:])
		repo.Save(uint32(i), data)
	}

	piece := make([]byte, lastPieceSize)
	for k := uint32(0); k < uint32(lastPieceSize); k++ {
		piece[k] = byte(k)
	}
	pieceHash := sha1.Sum(piece)
	pieceHashes = append(pieceHashes, pieceHash[:])

	repo.Save(uint32(numberOfPieces-1), piece)
	return repo, pieceHashes
}
