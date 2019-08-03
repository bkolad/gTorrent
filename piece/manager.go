package piece

import (
	"sync"

	"github.com/bkolad/gTorrent/torrent"
)

type Manager interface {
	SetNext([]bool, string) (bool, uint32)
	SetDone(int, []byte)
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
	sync.Mutex
	pieces []pieceStatus
	info   torrent.Info
}

func NewManager() Manager {
	return &manager{}
}

func (m *manager) SetNext(bitSet []bool, peerID string) (bool, uint32) {
	m.Lock()
	defer m.Unlock()

	for i, v := range bitSet {
		if v && m.pieces[i].empty() {
			m.pieces[i] = pieceStatus{peerID: peerID, done: false}
			return false, uint32(i)
		}
	}
	return true, 0
}

func (m *manager) SetDone(i int, b []byte) {
	m.Lock()
	defer m.Unlock()

	m.pieces[i] = pieceStatus{
		peerID: m.pieces[i].peerID,
		done:   true,
	}
}
