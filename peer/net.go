package peer

import (
	"io"
	"net"
	"strconv"
	"sync"
	"time"

	log "github.com/bkolad/gTorrent/logger"
	"github.com/bkolad/gTorrent/torrent"
)

// Network is the interface that has to be implemented
// in order to communicate with other bittorent peers,
// the concrete implementation has to use TCP as transport protocol.
type Network interface {
	// SendHandshake sends handshake to remote peer
	// and if successful registers connection for further communication.
	// SendHandshake should be called only once for every peer.
	SendHandshake() error
	// RegisterListener stores the listener. Only one listener can be registered.
	RegisterListener(Listener)
	// Send sends message to remote peer
	Send(Packet) error
}

// Listener defines a callback which will be invoked on every
// incoming packet.
type Listener interface {
	NewPacket(Packet)
}

type network struct {
	sync.Mutex
	peerInfo torrent.PeerInfo
	hanshake Handshake
	listener Listener
	conn     net.Conn
}

const dialerTimeOut = 10 * time.Second

// NewNetwork create Network
func NewNetwork(peerInfo torrent.PeerInfo, handshake Handshake) Network {
	return &network{
		peerInfo: peerInfo,
		hanshake: handshake,
	}
}

func (n *network) SendHandshake() error {
	addr := net.JoinHostPort(n.peerInfo.IP, strconv.Itoa(int(n.peerInfo.Port)))
	dialer := net.Dialer{Timeout: dialerTimeOut}
	conn, err := dialer.Dial("tcp", addr)

	if err != nil {
		log.Debug("Can't dial remote peer " + err.Error())
		return err
	}

	handshake, err := n.hanshake.Encode()
	if err != nil {
		log.Debug("Problem with handshake encoding")
		return err
	}
	_, err = conn.Write(handshake)
	if err != nil {
		log.Debug("Problem with sending handshake to " + addr + " " + err.Error())
		return err
	}
	//handshake takes 68 bytes
	var buf [68]byte

	_, err = io.ReadFull(conn, buf[:])
	if err != nil {
		log.Debug("Problem with receiving handshake from " + addr + " " + err.Error())
		return err
	}

	var handshakeFromPeer Handshake
	log.Info("Successfully received handshake from " + addr)

	err = handshakeFromPeer.Decode(buf[:])
	if err != nil {
		log.Debug("Problem with decoding handshake from " + addr + " " + err.Error())
		return err
	}
	n.conn = conn

	go func() {
		for {
			err := n.handleConn()
			if err != nil {
				return
			}
		}
	}()
	return nil
}

func (n *network) handleConn() error {
	packet := &packet{}
	err := packet.Decode(n.conn)
	if err != nil {
		log.Info("Packet from " + n.peerInfo.IP + " can't be decoded " + err.Error())
		return err
	}
	n.dispatch(packet)
	return nil
}

func (n *network) dispatch(p Packet) {
	n.Lock()
	defer n.Unlock()
	n.listener.NewPacket(p)
}

func (n *network) Send(p Packet) error {
	_, err := n.conn.Write(p.Encode())
	if err != nil {
		log.Error(err.Error())
	}
	return err
}

func (n *network) RegisterListener(l Listener) {
	n.listener = l
}
