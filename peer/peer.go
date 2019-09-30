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
	onPiece(uint32, uint32, []byte)
	onCancel()
	onPort()
	onUnknown()
	stop()
}

type simplePeer struct {
	msgs             chan MSG
	net              Network
	chocked          bool
	interested       bool
	pieceManager     p.Manager
	currentPiece     uint32
	currentPieceData []byte
	currentOffset    uint32
	peerInfo         torrent.PeerInfo
	done             bool
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
		msgs:         messages,
		net:          net,
		pieceManager: pieceManager,
		peerInfo:     peerInfo,
		done:         false,
	}
	net.RegisterListener(peer)
	return peer
}

func (p *simplePeer) start() {
	err := p.net.SendHandshake()
	if err != nil {
		log.Error(err.Error())
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
	done, next := p.pieceManager.NextPiece(p.peerInfo.IP)
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
	//data, err := p.pieceRepository.Get(piece, offset, size)

	pieceData, err := p.pieceManager.Get(piece)
	pieceChunk := pieceData[offset : offset+size]
	if err != nil {
		panic(err)
	}
	packet := encodePieceData(piece, offset, pieceChunk)
	p.send(packet)
}

func (p *simplePeer) onPiece(piece, offset uint32, payload []byte) {
	if p.currentOffset != offset {
		log.Error("Received bad offset " + p.peerInfo.IP)
		p.stop()
		return
	}
	p.currentOffset += uint32(len(payload))
	p.currentPieceData = append(p.currentPieceData, payload...)

	if p.pieceManager.IsLastChunk(piece, p.currentOffset) {
		err := p.pieceManager.PieceDone(piece, p.currentPieceData)
		if err != nil {
			log.Error("Cant save piece from " + p.peerInfo.IP)
			log.Error(err.Error())
			p.stop()
			return
		}

		done, nextPiece := p.pieceManager.NextPiece(p.peerInfo.IP)
		log.Info(p.peerInfo.IP + ": Downloaded piece: " + fmt.Sprint(p.currentPiece))
		if done {
			log.Info(p.peerInfo.IP + ": Downloading finisched")
			p.stop()
			return
		}
		p.currentPiece = nextPiece
		p.currentPieceData = make([]byte, 0)
		p.currentOffset = 0
	}
	packet := encodePieceRequest(p.currentPiece, p.currentOffset, p.pieceManager.ChunkSize())
	p.send(packet)
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

func (p *simplePeer) NewPacket(packet Packet) bool {
	if p.done {
		return true
	}

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
	return false
}

func (p *simplePeer) stop() {
	go func() {
		p.msgs <- kill{p.peerInfo}
	}()
	p.done = true
}
