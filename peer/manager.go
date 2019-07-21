package peer

import (
	log "github.com/bkolad/gTorrent/logger"
	"github.com/bkolad/gTorrent/torrent"
)

const maxActivePeers = 10

type Manager interface {
	ConnectToPeers()
}

type manager struct {
	peersInfo   chan torrent.PeerInfo
	activePeers map[torrent.PeerInfo]controller
	messages    chan MSG
	handshake   Handshake
}

func NewManager(peerInfoChan chan torrent.PeerInfo, handshake Handshake) Manager {
	activePeers := make(map[torrent.PeerInfo]controller)
	messages := make(chan MSG, 100)
	return &manager{peerInfoChan, activePeers, messages, handshake}
}

func (m *manager) ConnectToPeers() {
	for {
		for p := range m.peersInfo {
			if len(m.activePeers) < maxActivePeers {
				log.Info("connecting to peer " + p.IP)
				peerController := newController(m.messages, p, m.handshake)
				go peerController.start()
				m.activePeers[p] = peerController
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
