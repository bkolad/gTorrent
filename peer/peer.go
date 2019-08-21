package peer

import (
	"fmt"

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
	onPiece(uint32, uint32, []byte) bool
	onCancel()
	onPort()
	onUnknown()
}

type simplePeer struct {
	msgs             chan MSG
	net              Network
	chocked          bool
	interested       bool
	pieceManager     p.Manager
	pieceRepository  p.Repository
	currentPiece     uint32
	currentPieceData []byte
	currentOffset    uint32
	peerInfo         torrent.PeerInfo
}

func newPeer(messages chan MSG,
	peerInfo torrent.PeerInfo,
	handshake Handshake,
	pieceManager p.Manager,
) Peer {
	net := NewNetwork(peerInfo, handshake)
	return newPeerWithNetwork(net, messages, peerInfo, handshake, pieceManager)
}

func newPeerWithNetwork(net Network,
	messages chan MSG,
	peerInfo torrent.PeerInfo,
	handshake Handshake,
	pieceManager p.Manager,
) *simplePeer {
	peer := &simplePeer{
		msgs:            messages,
		net:             net,
		pieceManager:    pieceManager,
		peerInfo:        peerInfo,
		pieceRepository: p.NewRepo(pieceManager.PieceCount()),
	}
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
	done, next := p.pieceManager.SetNext(p.peerInfo.IP)
	if done {
		return
	}
	p.currentPiece = next
	packet := encodePieceRequest(next, 0, p.pieceManager.ChunkSize())
	p.currentOffset = 0
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
	//idx := haveToIndex(payload)
	//p.remotePeerPieces[idx] = true
	//packet := encodeInterested()
	//p.send(packet)
}

func (p *simplePeer) onBitfield(bitfield []byte) {
	remotePeerPieces := bytesToBits(bitfield)
	p.pieceManager.SetPeerPieces(p.peerInfo.IP, remotePeerPieces)
	packet := encodeInterested()
	p.send(packet)
}

func (p *simplePeer) onRequest(piece, offset, size uint32) {
	data := p.pieceRepository.Get(piece, offset, size)
	packet := encodePieceData(piece, offset, data)
	p.send(packet)
}

func (p *simplePeer) onPiece(piece, offset uint32, payload []byte) bool {
	p.currentOffset += uint32(len(payload))
	p.currentPieceData = append(p.currentPieceData, payload...)

	isLastChunk := p.currentOffset == p.pieceManager.PieceSize(piece)
	if isLastChunk {
		p.pieceRepository.Save(piece, p.currentPieceData)
		done, nextPiece := p.pieceManager.SetNext(p.peerInfo.IP)
		log.Info(p.peerInfo.IP + ": Downloaded piece: " + fmt.Sprint(p.currentPiece))
		if done {
			return true
		}
		p.currentPiece = nextPiece
		p.currentPieceData = make([]byte, 0)
		p.currentOffset = 0
	}
	packet := encodePieceRequest(p.currentPiece, p.currentOffset, p.pieceManager.ChunkSize())
	p.send(packet)
	return false
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
