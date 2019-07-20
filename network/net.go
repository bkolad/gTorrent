package network

import (
	"io"
	"net"
	"strconv"
	"time"

	"github.com/bkolad/gTorrent/peer"
	"github.com/bkolad/gTorrent/torrent"
)

type Network interface {
	Send()
	RegisterListener(Listener)
}

type Listener interface {
	NewPacket()
}

type network struct {
	peerInfo *torrent.PeerInfo
	h        *peer.Handshake
}

const dialerTimeOut = 10 * time.Second

func NewNetwork(peerInfo *torrent.PeerInfo, h *peer.Handshake) Network {
	return &network{peerInfo, h}
}

func (n *network) Send() {
	addr := net.JoinHostPort(n.peerInfo.IP, strconv.Itoa(int(n.peerInfo.Port)))
	dialer := net.Dialer{Timeout: dialerTimeOut}
	conn, err := dialer.Dial("tcp", addr)

	if err != nil {
		panic(err)
	}

	handshake, err := n.h.Encode()
	if err != nil {
		panic(err)
	}

	conn.Write(handshake)
	var buf [68]byte
	_, _ = io.ReadFull(conn, buf[:])

	//handshakeFromPeer, err := peer.Decode(buf[:])
	//if err != nil {
	//	panic(err)
	//}

	//log.Default.Info("handshake ")

}

func (n *network) RegisterListener(l Listener) {

}
