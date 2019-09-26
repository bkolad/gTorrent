package piece

import (
	"bytes"
	"crypto/sha1"
	"sync"

	"github.com/bkolad/gTorrent/torrent"
)

type Manager interface {
	SetNext(string) (bool, uint32)
	SetDone(uint32, []byte)
	SetPeerPieces(string, []bool)
	PieceSize(n uint32) uint32
	ChunkSize() uint32
	PieceCount() uint32
	Verify(i uint32, piece []byte) bool
	//Done() []int
	//InProgress() []int
	//PieceLength() int
	//LastPieceLength() int
	//PieceHash(int) []byte
}

type pieceStatus struct {
	peerID string
	done   bool
}

func (p pieceStatus) empty() bool {
	return p.peerID == ""
}

func (p pieceStatus) have() bool {
	return !p.empty() && p.done
}

func (p pieceStatus) inProgress() bool {
	return !p.empty() && !p.done
}

type manager struct {
	info          torrent.Info
	lastPieceSize uint32
	sync.Mutex
	pieces      []pieceStatus
	peersPieces map[string][]bool
	pieceSize   uint32
	pieceRepo   Repository
}

func NewManager(info torrent.Info, pieceRepo Repository) *manager {
	lastPieceSize, pieceCount := torrent.CalculateLastPieceSize(info.Length, info.PieceSize)
	pieces := make([]pieceStatus, pieceCount)

	return &manager{
		info:          info,
		pieces:        pieces,
		peersPieces:   map[string][]bool{},
		lastPieceSize: uint32(lastPieceSize),
		pieceSize:     uint32(info.PieceSize),
		pieceRepo:     pieceRepo,
	}
}

func (m *manager) PieceCount() uint32 {
	return uint32(len(m.pieces))
}

func (m *manager) PieceSize(n uint32) uint32 {
	if n == m.PieceCount()-1 {
		return m.lastPieceSize
	}
	return m.pieceSize
}

func (m *manager) ChunkSize() uint32 {
	return uint32(m.info.ChunkSize)
}

func (m *manager) SetNext(peerID string) (bool, uint32) {
	m.Lock()
	defer m.Unlock()

	peerPieces := m.peersPieces[peerID]
	for i, v := range peerPieces {
		if v && m.pieces[i].empty() {
			m.pieces[i] = pieceStatus{peerID: peerID, done: false}
			return false, uint32(i)
		}
	}
	return true, 0
}

func (m *manager) SetDone(i uint32, b []byte) {
	m.Lock()
	defer m.Unlock()
	m.pieces[i] = pieceStatus{
		peerID: m.pieces[i].peerID,
		done:   true,
	}
}

func (m *manager) SetPeerPieces(peerID string, pieces []bool) {
	m.Lock()
	defer m.Unlock()
	m.peersPieces[peerID] = pieces
}

func (m *manager) Verify(i uint32, piece []byte) bool {
	pieceHash := sha1.Sum(piece)
	return 0 == bytes.Compare(m.info.PieceHashes[i], pieceHash[:])
}
