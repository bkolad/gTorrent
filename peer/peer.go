package peer

import (
	"fmt"
	"strconv"

	log "github.com/bkolad/gTorrent/logger"
	"github.com/bkolad/gTorrent/torrent"
)

type Peer interface {
	start()
	onChoke()
	onUnchoke()
	onInterested()
	onNotInterested()
	onHave([]byte)
	onBitfield([]byte)
	onRequest()
	onPiece(uint32, uint32, []byte)
	onCancel()
	onPort()
	onUnknown()
}

type simplePeer struct {
	msgs       chan MSG
	net        Network
	chocked    bool
	bitfield   []byte
	interested bool
}

func newPeer(messages chan MSG, peerInfo torrent.PeerInfo, handshake Handshake) Peer {
	net := NewNetwork(peerInfo, handshake)
	peer := &simplePeer{msgs: messages, net: net}
	net.RegisterListener(peer)
	return peer
}

func (p *simplePeer) start() {
	err := p.net.SendHandshake()
	if err != nil {
		fmt.Println("Err" + err.Error())
		p.msgs <- handshakeError{}
	}
}

func (p *simplePeer) onKeepAlive() {
}

func (p *simplePeer) onChoke() {
	p.chocked = true
}

func (p *simplePeer) onUnchoke() {
	log.Debug("Unchoked")
	p.chocked = false
	packet := encodePieceRequest(0, 0, 16384)
	p.send(packet)
}

func (p *simplePeer) onInterested() {
	p.interested = true
}

func (p *simplePeer) onNotInterested() {
	p.interested = false
}

func (p *simplePeer) onHave(payload []byte) {
	log.Debug("have")
	packet := encodeInterested()
	p.send(packet)
}

func (p *simplePeer) onBitfield(bitfield []byte) {
	p.bitfield = bitfield
	packet := encodeInterested()
	p.send(packet)
}

func (p *simplePeer) onRequest() {

}

func (p *simplePeer) onPiece(piece, offset uint32, payload []byte) {
	log.Debug("Received piece " + strconv.Itoa(int(piece)) + "  " + strconv.Itoa(int(offset)) + " " + strconv.Itoa(len(payload)))
}

func (p *simplePeer) onCancel() {

}

func (p *simplePeer) onPort() {

}

func (p *simplePeer) onUnknown() {

}

func (p *simplePeer) send(packet Packet) {
	p.net.Send(packet)
}

func (p *simplePeer) NewPacket(packet Packet) {
	switch int(packet.ID()) {
	case keepAlaive:
		p.onKeepAlive()
	case choke:
		p.onChoke()
	case unchoke:
		p.onUnchoke()
	case interested:
		p.onInterested()
	case notInterested:
		p.onNotInterested()
	case have:
		p.onHave(packet.Payload())
	case bitfield:
		p.onBitfield(packet.Payload())
	case request:
		p.onRequest()
	case piece:
		piece, offset, data := decodePiece(packet.Payload())
		p.onPiece(piece, offset, data)
	case cancel:
		p.onCancel()
	case port:
		p.onPort()
	case unknown:
		p.onUnknown()
	}
}
