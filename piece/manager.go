package piece

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"strconv"
	"sync"

	"github.com/bkolad/gTorrent/torrent"
)

//Manager manages file pieces, in bittorent protocol file is divided into pieces.
//Torrent file contains information about piece size and piece hash.
//Implementation of this interface must be thread safe.
type Manager interface {
	//SetPeerPieces associates remote peer and pieces it possess
	SetPeerPieces(peerID string, pieces []bool)
	//PieceSize returns size of a piece corresponding to given index,
	//often last piece has different size than preceding pieces.
	PieceSize(index uint32) uint32
	ChunkSize() uint32
	//Get retrieves piece corresponding to given index
	Get(index uint32) (pieceData []byte, err error)
	//PieceDone saves piece in the storage
	PieceDone(index uint32, pieceData []byte) error
	//NextPiece checks if peers has more pieces to offer.
	// (false, 0) -> peer has no more pieces
	// (true, index) -> index is the next piece index to request
	NextPiece(peerID string) (hasNext bool, index uint32)
	//IsLastChunk checks if the whole piece was downloaded
	IsLastChunk(pieceIndex, offset uint32) bool
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

//NewManager creates new piece manager
func NewManager(info torrent.Info, pieceRepo Repository) *manager {
	lastPieceSize, pieceCount := info.CalculateLastPieceSize()
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

func (m *manager) SetPeerPieces(peerID string, pieces []bool) {
	m.Lock()
	defer m.Unlock()
	m.peersPieces[peerID] = pieces
}

func (m *manager) PieceSize(n uint32) uint32 {
	if n == uint32(len(m.pieces))-1 {
		return m.lastPieceSize
	}
	return m.pieceSize
}

func (m *manager) ChunkSize() uint32 {
	return uint32(m.info.ChunkSize)
}

func (m *manager) NextPiece(peerID string) (bool, uint32) {
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

func (m *manager) verify(index uint32, pieceData []byte) bool {
	pieceHash := sha1.Sum(pieceData)
	return 0 == bytes.Compare(m.info.PieceHashes[index], pieceHash[:])
}

func (m *manager) save(piece uint32, data []byte) error {
	return m.pieceRepo.Save(piece, data)
}

func (m *manager) Get(piece uint32) ([]byte, error) {
	return m.pieceRepo.Get(piece)
}

func (m *manager) PieceDone(
	index uint32,
	pieceData []byte) error {

	ok := m.verify(index, pieceData)
	if !ok {
		return errors.New("Wrong hash for piece " + strconv.Itoa(int(index)))
	}

	err := m.pieceRepo.Save(index, pieceData)
	if err != nil {
		return err
	}

	m.Lock()
	m.pieces[index] = pieceStatus{
		peerID: m.pieces[index].peerID,
		done:   true,
	}

	m.Unlock()
	return nil
}

func (m *manager) IsLastChunk(pieceIndex, offset uint32) bool {
	return offset == m.PieceSize(pieceIndex)
}
