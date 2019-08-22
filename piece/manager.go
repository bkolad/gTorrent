package piece

import (
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
	return p.peerID != "" && p.done
}

func (p pieceStatus) inProgress() bool {
	return p.peerID != "" && !p.done
}

type manager struct {
	info          torrent.Info
	lastPieceSize uint32
	sync.Mutex
	pieces      []pieceStatus
	peersPieces map[string][]bool
	pieceSize   uint32
}

func NewManager(info torrent.Info) *manager {
	lastPieceSize, numberOfPieces := CalculateLastPieceSize(info.Length, info.PieceSize)

	pieces := make([]pieceStatus, numberOfPieces)
	return &manager{
		info:          info,
		pieces:        pieces,
		peersPieces:   map[string][]bool{},
		lastPieceSize: uint32(lastPieceSize),
		pieceSize:     uint32(info.PieceSize),
	}
}

func CalculateLastPieceSize(length int, pieceSize int) (uint32, uint32) {
	lastPieceSize := length % pieceSize
	numberOfPieces := length / pieceSize

	if lastPieceSize != 0 {
		numberOfPieces++
	} else {
		lastPieceSize = pieceSize
	}
	return uint32(lastPieceSize), uint32(numberOfPieces)
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
