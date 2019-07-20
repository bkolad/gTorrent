package peer

import (
	"fmt"

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
	c2 := make(map[torrent.PeerInfo]controller)
	c3 := make(chan MSG, 100)
	return &manager{peerInfoChan, c2, c3, handshake}
}

func (m *manager) ConnectToPeers() {
	fmt.Println("ConnectToPeers ")

	for {
		select {
		case p := <-m.peersInfo:
			fmt.Println("lol ", p)

			for len(m.activePeers) < maxActivePeers {
				fmt.Println("add ", p)
				peerController := newController(m.messages, p, m.handshake)
				go peerController.start()
				m.activePeers[p] = peerController
			}
			break

		}

		for msg := range m.messages {
			switch msg := msg.(type) {
			case killed:
				close(m.activePeers[msg.peerInfo].c)
				delete(m.activePeers, msg.peerInfo)
			} //TODO statistics
		}
	}
}
