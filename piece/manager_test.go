package piece

import (
	"testing"

	"github.com/bkolad/gTorrent/torrent"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	info           torrent.Info
	numberOfPieces int
	lastPieceSize  uint32
}

func tests() []testCase {
	return []testCase{
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
			lastPieceSize:  80,
		},
		testCase{
			info: torrent.Info{
				PieceSize: 120,
				Length:    120},
			numberOfPieces: 1,
			lastPieceSize:  120,
		},
	}
}

func TestPeerManager(t *testing.T) {
	for _, testCase := range tests() {
		manager := NewManager(testCase.info, nil)
		require.Equal(t, int(testCase.lastPieceSize), int(manager.lastPieceSize))
		require.Equal(t, testCase.numberOfPieces, len(manager.pieces))
	}
}

func TestPeerManagerNexPiece(t *testing.T) {
	for _, testCase := range tests() {
		missing := 1
		peer1 := makePeer("peer1", testCase, 1, missing)
		peer2 := makePeer("peer2", testCase, 2, missing)
		peer3 := makePeer("peer3", testCase, 3, missing)
		repo := NewRepo(uint32(testCase.numberOfPieces))

		manager := NewManager(testCase.info, repo)
		manager.SetPeerPieces(peer1.peerID, peer1.pieces)
		manager.SetPeerPieces(peer2.peerID, peer2.pieces)
		manager.SetPeerPieces(peer3.peerID, peer3.pieces)

		done := false
		for i := 0; i < testCase.numberOfPieces; i++ {
			done, next := manager.NextPiece(peer1.peerID)
			if !done {
				require.True(t, peer1.pieces[next])
			}

			done, next = manager.NextPiece(peer2.peerID)
			if !done {
				require.True(t, peer2.pieces[next])
			}

			done, next = manager.NextPiece(peer3.peerID)
			if !done {
				require.True(t, peer3.pieces[next])
			}
		}
		done, _ = manager.NextPiece(peer1.peerID)
		require.True(t, done)
		done, _ = manager.NextPiece(peer2.peerID)
		require.True(t, done)
		done, _ = manager.NextPiece(peer3.peerID)
		require.True(t, done)

		if missing < testCase.numberOfPieces {
			require.True(t, manager.pieces[missing].empty())
		}
	}
}

type peer struct {
	peerID string
	pieces []bool
}

func makePeer(peerID string, t testCase, n int, missing int) peer {
	pieces := make([]bool, t.numberOfPieces)

	for i := 0; i < len(pieces); i++ {
		if i%n == 0 && i != missing {
			pieces[i] = true
		}
	}
	return peer{peerID, pieces}
}
