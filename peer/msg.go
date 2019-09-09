package peer

import "github.com/bkolad/gTorrent/torrent"

type MSG interface {
}

type handshakeError struct{}

type kill struct {
	peerInfo torrent.PeerInfo
}
