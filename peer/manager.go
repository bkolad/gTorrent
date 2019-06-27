package peer

import "github.com/bkolad/gTorrent/tracker"

type Manager interface {
	ActivePeers() []tracker.PeerInfo
}
