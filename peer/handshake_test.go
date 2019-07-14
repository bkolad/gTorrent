package peer

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHandshake(t *testing.T) {
	handshake := Handshake{
		len:      19,
		protocol: "BitTorrent protocol",
		rsvd:     make([]byte, 8),
		InfoHash: []byte("sdfvcxujklppmntrqzxc"),
		peerID:   "-GT001-" + "d0p22uiake0bd",
	}

	data, err := handshake.Encode()
	require.NoError(t, err)
	decodedHandshake, err := Decode(data)
	require.NoError(t, err)
	require.Equal(t, handshake, decodedHandshake)
}
