package peer

import (
	"io"
	"net"
	"strconv"
	"time"

	log "github.com/bkolad/gTorrent/logger"
	"github.com/bkolad/gTorrent/torrent"
)

type Network interface {
	SendHandshake() error
	RegisterListener(Listener)
}

type Listener interface {
	NewPacket()
}

type network struct {
	peerInfo torrent.PeerInfo
	hanshake Handshake
}

const dialerTimeOut = 20 * time.Second

func NewNetwork(peerInfo torrent.PeerInfo, handshake Handshake) Network {
	return &network{peerInfo, handshake}
}

func (n *network) SendHandshake() error {
	addr := net.JoinHostPort(n.peerInfo.IP, strconv.Itoa(int(n.peerInfo.Port)))
	dialer := net.Dialer{Timeout: dialerTimeOut}
	conn, err := dialer.Dial("tcp", addr)

	if err != nil {
		panic(err)
	}

	handshake, err := n.hanshake.Encode()
	if err != nil {
		return err
	}

	conn.Write(handshake)
	var buf [68]byte

	_, _ = io.ReadFull(conn, buf[:])

	var handshakeFromPeer Handshake
	err = handshakeFromPeer.Decode(buf[:])
	if err != nil {
		return err
	}

	log.Info("handshake " + addr)
	return nil
}

func (n *network) RegisterListener(l Listener) {

}
