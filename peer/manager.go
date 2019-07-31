package peer

import (
	log "github.com/bkolad/gTorrent/logger"
	"github.com/bkolad/gTorrent/torrent"
)

const maxActivePeers = 30

type Manager interface {
	ConnectToPeers()
}

type manager struct {
	peersInfo   chan torrent.PeerInfo
	activePeers map[torrent.PeerInfo]Peer
	messages    chan MSG
	handshake   Handshake
}

func NewManager(peerInfoChan chan torrent.PeerInfo, handshake Handshake) Manager {
	activePeers := make(map[torrent.PeerInfo]Peer)
	messages := make(chan MSG, 100)
	return &manager{peerInfoChan, activePeers, messages, handshake}
}

func (m *manager) ConnectToPeers() {
	for {
		for p := range m.peersInfo {
			if len(m.activePeers) < maxActivePeers {
				log.Info("connecting to peer " + p.IP)
				peer := newPeer(m.messages, p, m.handshake)
				go peer.start()
				m.activePeers[p] = peer
			} else {
				break
			}
		}

		for msg := range m.messages {
			switch msg := msg.(type) {
			case killed:
				//	close
				delete(m.activePeers, msg.peerInfo)
			case handshakeError:
				log.Error("HandshakeError")
			}
		}
	}
}
