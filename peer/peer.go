package peer

import "fmt"
import "github.com/bkolad/gTorrent/torrent"

type Peer interface {
	start()
	onChoke()
	onUnchoke()
	onBitfield()
	send()
}

type simplePeer struct {
	msgs chan MSG
	net  Network
}

func newPeer(messages chan MSG, peerInfo torrent.PeerInfo, handshake Handshake) Peer {
	net := NewNetwork(peerInfo, handshake)
	peer := simplePeer{msgs: messages, net: net}
	net.RegisterListener(peer)
	return &peer
}

func (p *simplePeer) start() {
	err := p.net.SendHandshake()
	if err != nil {
		fmt.Println("Err" + err.Error())
		p.msgs <- handshakeError{}
	}
}

func (p *simplePeer) onChoke() {
	fmt.Println("choke")
}

func (p *simplePeer) onUnchoke() {
	fmt.Println("unchoke")
}

func (p *simplePeer) onBitfield() {
	fmt.Println("BF")
}

func (p *simplePeer) send() {
	fmt.Println("send")
}

func (p simplePeer) NewPacket(packet Packet) {
	fmt.Println("New Packet", packet.ID())
}
