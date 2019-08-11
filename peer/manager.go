package peer

import (
	"fmt"

	log "github.com/bkolad/gTorrent/logger"
	p "github.com/bkolad/gTorrent/piece"
	"github.com/bkolad/gTorrent/torrent"
)

const maxActivePeers = 30

type Manager interface {
	ConnectToPeers()
}

type manager struct {
	peersInfo    chan torrent.PeerInfo
	activePeers  map[torrent.PeerInfo]Peer
	messages     chan MSG
	handshake    Handshake
	pieceManager p.Manager
}

func NewManager(peersInfo chan torrent.PeerInfo,
	handshake Handshake,
	pieceManager p.Manager,
) Manager {
	activePeers := make(map[torrent.PeerInfo]Peer)
	messages := make(chan MSG, 100)
	return &manager{
		peersInfo:    peersInfo,
		activePeers:  activePeers,
		messages:     messages,
		handshake:    handshake,
		pieceManager: pieceManager,
	}
}

func (m *manager) ConnectToPeers() {
	for {
		for p := range m.peersInfo {
			if len(m.activePeers) < maxActivePeers {
				log.Info("connecting to peer " + p.IP)
				peer := newPeer(m.messages, p, m.handshake, m.pieceManager)
				go peer.start()
				//TODO fix reace condition
				m.activePeers[p] = peer
				fmt.Println(m.activePeers[p])
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
