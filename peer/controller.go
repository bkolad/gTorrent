package peer

import "github.com/bkolad/gTorrent/torrent"

type controller struct {
	messages chan MSG
	net      Network
}

func newController(messages chan MSG, peerInfo torrent.PeerInfo, handshake Handshake) controller {
	return controller{messages, NewNetwork(peerInfo, handshake)}
}

func (c controller) start() {
	peer := peer{c.messages, c.net}
	peer.doHandshake()
}
