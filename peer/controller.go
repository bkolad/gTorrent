package peer

import "github.com/bkolad/gTorrent/torrent"

type controller struct {
	c   chan MSG
	net Network
}

func newController(msgs chan MSG, peerInfo torrent.PeerInfo, handshake Handshake) controller {

	return controller{nil, NewNetwork(&peerInfo, handshake)}
}

func (c controller) start() {
	peer := peer{nil, c.net}
	peer.doHandshake()
}
