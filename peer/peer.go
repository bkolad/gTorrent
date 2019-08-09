package peer

import (
	"fmt"
	"strconv"

	log "github.com/bkolad/gTorrent/logger"
	p "github.com/bkolad/gTorrent/piece"
	"github.com/bkolad/gTorrent/torrent"
)

// Peer handels messages form other bittorent nodes.
// - requests new pieces
// - serves pieces which are available to the local peer
// - tracks who has what content
// ... see BitTorrent spec: https://www.bittorrent.org/beps/bep_0003.html
type Peer interface {
	start()
	onChoke()
	onUnchoke()
	onInterested()
	onNotInterested()
	onHave([]byte)
	onBitfield([]byte)
	onRequest(uint32, uint32, uint32)
	onPiece(uint32, uint32, []byte)
	onCancel()
	onPort()
	onUnknown()
}

type simplePeer struct {
	msgs             chan MSG
	net              Network
	chocked          bool
	remotePeerPieces []bool
	interested       bool
	pieceManager     p.Manager
	currentPiece     []byte
	currentOffset    int
}

func newPeer(messages chan MSG,
	peerInfo torrent.PeerInfo,
	handshake Handshake,
	pieceManager p.Manager,
) Peer {
	net := NewNetwork(peerInfo, handshake)
	peer := &simplePeer{msgs: messages, net: net, pieceManager: pieceManager}
	net.RegisterListener(peer)
	return peer
}

func (p *simplePeer) start() {
	err := p.net.SendHandshake()
	if err != nil {
		fmt.Println("Err " + err.Error())
		p.msgs <- handshakeError{}
	}
}

// callbacks runs allways in the same "peer" go-routine

func (p *simplePeer) onKeepAlive() {
	log.Debug("keep alive")
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
	idx := haveToIndex(payload)
	p.remotePeerPieces[idx] = true
	packet := encodeInterested()
	p.send(packet)
}

func (p *simplePeer) onBitfield(bitfield []byte) {
	p.remotePeerPieces = bytesToBits(bitfield)
	packet := encodeInterested()
	p.send(packet)
}

func (p *simplePeer) onRequest(piece, offset, size uint32) {

}

func (p *simplePeer) onPiece(piece, offset uint32, payload []byte) {
	//	p.current = append(p.current, payload)
	p.currentOffset += len(payload)
	lastBlock := false
	if lastBlock {
		// done, nextPiece := p.pieceManager.SetNext(bitSet, "peerId")
		// p.currentPiece = make([]byte, 0)
		//packet := encodePieceRequest(nextPiece, 0, 16384)
		// send
	}
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
	switch packet.ID() {
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
		p.onRequest(decodeRequest(packet.Payload()))
	case piece:
		p.onPiece(decodePiece(packet.Payload()))
	case cancel:
		p.onCancel()
	case port:
		p.onPort()
	case unknown:
		p.onUnknown()
	}
}
