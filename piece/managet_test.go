package piece

import (
	"testing"

	"github.com/bkolad/gTorrent/torrent"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	info           torrent.Info
	numberOfPieces int
	lastPieceSize  int
}

func TestPeerManager(t *testing.T) {
	testCases := []testCase{
		testCase{
			info: torrent.Info{
				PieceSize: 21,
				Length:    100},
			numberOfPieces: 5,
			lastPieceSize:  16,
		},
		testCase{
			info: torrent.Info{
				PieceSize: 8,
				Length:    80},
			numberOfPieces: 10,
			lastPieceSize:  8,
		},
		testCase{
			info: torrent.Info{
				PieceSize: 120,
				Length:    80},
			numberOfPieces: 1,
			lastPieceSize:  120,
		},
	}

	for _, testCase := range testCases {
		manager := NewManager(testCase.info)
		require.Equal(t, testCase.lastPieceSize, manager.lastPieceSize)
		require.Equal(t, testCase.numberOfPieces, len(manager.pieces))
	}
}

func TestPeerManagerNexPiece(t *testing.T) {

}
