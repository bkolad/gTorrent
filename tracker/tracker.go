package tracker

import "github.com/bkolad/gTorrent/torrent"

//Tracker retrieves peers from the tracker
type Tracker interface {
	Peers() ([]*torrent.PeerInfo, error)
}
