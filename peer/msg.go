package peer

import "github.com/bkolad/gTorrent/torrent"

type MSG interface {
}

type killed struct {
	peerInfo torrent.PeerInfo
}
