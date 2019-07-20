package peer

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/bkolad/gTorrent/torrent"
)

type Network interface {
	SendHandshake()
	RegisterListener(Listener)
}

type Listener interface {
	NewPacket()
}

type network struct {
	peerInfo *torrent.PeerInfo
	h        Handshake
}

const dialerTimeOut = 20 * time.Second

func NewNetwork(peerInfo *torrent.PeerInfo, h Handshake) Network {
	return &network{peerInfo, h}
}

func (n *network) SendHandshake() {
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

	var handshakeFromPeer Handshake
	err = handshakeFromPeer.Decode(buf[:])
	if err != nil {
		panic(err)
	}

	fmt.Println("handshake ", handshakeFromPeer)

}

func (n *network) RegisterListener(l Listener) {

}
